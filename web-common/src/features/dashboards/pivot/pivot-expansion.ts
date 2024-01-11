import type { StateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";
import {
  createPivotAggregationRowQuery,
  getAxisForDimensions,
} from "./pivot-data-store";
import type { ExpandedState } from "@tanstack/svelte-table";
import { derived, writable } from "svelte/store";
import type { PivotDataStoreConfig } from "./types";
import { getFilterForPivotTable } from "./pivot-utils";
import {
  createTableWithAxes,
  reduceTableCellDataIntoRows,
} from "./pivot-table-transformations";

/**
 * Extracts and organizes dimension names from a nested array structure
 * based on a specified anchor dimension and an expanded state.
 *
 * This function iterates over each key in the `expanded` object, which
 * indicates whether a particular path in the nested array is expanded.
 * For each expanded path, it navigates through the table data
 * following the path defined by the key (split into indices) and extracts
 * the dimension values at each level.
 *
 */
function getExpandedValuesFromNestedArray(
  tableData,
  anchorDimension: string,
  expanded: ExpandedState,
): Record<string, string[]> {
  const values = {};

  for (const key in expanded as Record<string, boolean>) {
    if (expanded[key]) {
      // Split the key into indices
      const indices = key.split(".").map((index) => parseInt(index, 10));

      // Retrieve the value from the nested array
      let currentValue = tableData;
      const dimensionNames: string[] = [];
      for (const index of indices) {
        if (!currentValue?.[index]) break;
        dimensionNames.push(currentValue[index]?.[anchorDimension]);
        currentValue = currentValue[index]?.subRows;
      }

      // Add the value to the result array
      values[key] = dimensionNames;
    }
  }

  return values;
}

/**
 * Returns a query for cell data for a sub table. The values are
 * sorted by anchor dimension so that then can be reduced into
 * rows optimally.
 */
export function createSubTableCellQuery(
  ctx: StateManagers,
  config: PivotDataStoreConfig,
  anchorDimension: string,
  columnDimensionAxesData,
  rowNestFilters,
) {
  const allDimensions = config.colDimensionNames.concat([anchorDimension]);

  const filterForSubTable = getFilterForPivotTable(
    config,
    columnDimensionAxesData,
  );

  const includeFilters = filterForSubTable.include.concat(rowNestFilters);
  const filters = {
    include: includeFilters,
    exclude: [],
  };

  const sortBy = [
    {
      desc: false,
      name: anchorDimension,
    },
  ];
  return createPivotAggregationRowQuery(
    ctx,
    config.measureNames,
    allDimensions,
    filters,
    sortBy,
    "10000",
  );
}

/**
 * For each expanded row, create a query for the sub table
 * and return the query result along with the expanded row index
 * and the row dimension values
 */
export function queryExpandedRowMeasureValues(
  ctx: StateManagers,
  config: PivotDataStoreConfig,
  tableData,
  columnDimensionAxesData,
) {
  const { rowDimensionNames } = config;
  const expanded = config.pivot.expanded;
  if (!tableData || Object.keys(expanded).length == 0) return writable(null);
  const values = getExpandedValuesFromNestedArray(
    tableData,
    rowDimensionNames[0],
    expanded,
  );

  return derived(
    Object.keys(values)?.map((expandIndex) => {
      const anchorDimension = rowDimensionNames[values[expandIndex].length];
      // TODO: handle for already existing filters
      const rowNestFilters = values[expandIndex].map((value, index) => {
        return {
          name: rowDimensionNames[index],
          in: [value],
        };
      });

      const filterForRowDimensionAxes = {
        include: rowNestFilters,
        exclude: [],
      };

      return derived(
        [
          writable(expandIndex),
          getAxisForDimensions(
            ctx,
            [anchorDimension],
            filterForRowDimensionAxes,
          ),
          createSubTableCellQuery(
            ctx,
            config,
            anchorDimension,
            columnDimensionAxesData,
            rowNestFilters,
          ),
        ],
        ([expandIndex, subRowDimensionValues, subTableData]) => {
          return {
            isFetching: subTableData?.isFetching,
            expandIndex,
            rowDimensionValues: subRowDimensionValues?.data?.[anchorDimension],
            data: subTableData?.data?.data,
          };
        },
      );
    }),
    (combos) => {
      return combos;
    },
  );
}

/***
 * For each expanded row, add the sub table data to the pivot table
 * data at the correct position.
 *
 * Note: Since the nested dimension values are present in the outermost
 * dimension's column, their accessor is the same as the anchor dimension.
 * Therefore, we change the key of the nested dimension to the anchor.
 */
export function addExpandedDataToPivot(
  config: PivotDataStoreConfig,
  tableData: Array<{ [key: string]: unknown }>,
  rowDimensions: string[],
  columnDimensionAxes,
  expandedRowMeasureValues,
): Array<{ [key: string]: unknown }> {
  const pivotData = tableData;
  const numRowDimensions = rowDimensions.length;

  expandedRowMeasureValues.forEach((expandedRowData) => {
    const indices = expandedRowData.expandIndex
      .split(".")
      .map((index) => parseInt(index, 10));

    let parent = pivotData; // Keep a reference to the parent array
    let lastIdx = 0;

    // Traverse the data array to the right position
    for (let i = 0; i < indices.length; i++) {
      if (!parent[indices[i]]) break;
      if (i < indices.length - 1) {
        parent = parent[indices[i]].subRows;
      }
      lastIdx = indices[i];
    }

    // Update the specific array at the position
    if (parent[lastIdx] && parent[lastIdx].subRows) {
      const anchorDimension = rowDimensions[indices.length];
      const rowValues = expandedRowData.rowDimensionValues;

      let skeletonSubTable: Array<{ [key: string]: unknown }> = [
        { [anchorDimension]: "LOADING_CELL" },
      ];
      if (expandedRowData?.rowDimensionValues?.length) {
        skeletonSubTable = createTableWithAxes(anchorDimension, rowValues);
      }

      let subTableData = skeletonSubTable;
      if (expandedRowData?.data?.length) {
        subTableData = reduceTableCellDataIntoRows(
          config,
          anchorDimension,
          expandedRowData?.rowDimensionValues,
          columnDimensionAxes?.data,
          skeletonSubTable,
          expandedRowData?.data,
        );
      }

      parent[lastIdx].subRows = subTableData?.map((row) => {
        const newRow = {
          ...row,
          [rowDimensions[0]]: row[anchorDimension],
        };

        /**
         * Add sub rows to the new row if number of row dimensions
         * is greater than number of nest levels expanded except
         * for the last level
         */
        if (numRowDimensions - 1 > indices.length) {
          newRow.subRows = [{ [rowDimensions[0]]: "LOADING_CELL" }];
        }
        return newRow;
      });
    }
  });
  return pivotData;
}
