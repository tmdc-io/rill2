import type {
  MetricsViewFilterCond,
  MetricsViewSpecDimensionV2,
  MetricsViewSpecMeasureV2,
} from "@rilldata/web-common/runtime-client";
import PivotExpandableCell from "./PivotExpandableCell.svelte";
import type { PivotDataRow, PivotDataStoreConfig, PivotState } from "./types";
import type { ColumnDef } from "@tanstack/svelte-table";

export function getMeasuresInPivotColumns(
  pivot: PivotState,
  measures: MetricsViewSpecMeasureV2[],
): MetricsViewSpecMeasureV2[] {
  const { columns } = pivot;

  return columns
    .filter((rowName) => measures.findIndex((m) => m?.name === rowName) > -1)
    .map((rowName) => measures.find((m) => m?.name === rowName));
}

export function getDimensionsInPivotRow(
  pivot: PivotState,
  dimensions: MetricsViewSpecDimensionV2[],
): MetricsViewSpecDimensionV2[] {
  const { rows } = pivot;
  return rows
    .filter(
      (rowName) => dimensions.findIndex((m) => m?.column === rowName) > -1,
    )
    .map((rowName) => dimensions.find((m) => m?.column === rowName));
}

export function getDimensionsInPivotColumns(
  pivot: PivotState,
  dimensions: MetricsViewSpecDimensionV2[],
): MetricsViewSpecDimensionV2[] {
  const { columns } = pivot;
  return columns
    .filter(
      (colName) => dimensions.findIndex((m) => m?.column === colName) > -1,
    )
    .map((colName) => dimensions.find((m) => m?.column === colName));
}

/**
 * Alternative to flexRender for performant rendering of cells
 */
export const cellComponent = (
  component: unknown,
  props: Record<string, unknown>,
) => ({
  component,
  props,
});

/**
 * Create a value to index map for a given array
 */
export function createIndexMap<T>(arr: T[]): Map<T, number> {
  const indexMap = new Map<T, number>();
  arr.forEach((element, index) => {
    indexMap.set(element, index);
  });
  return indexMap;
}

/***
 * Get filter for table cells
 */
export function getFilterForPivotTable(
  config: PivotDataStoreConfig,
  colDimensionAxes: Record<string, string[]> = {},
  rowDimensionAxes: Record<string, string[]> = {},
  isInitialTable = false,
  yLimit = 100,
  xLimit = 20,
) {
  // TODO: handle for already existing global filters

  const { colDimensionNames, rowDimensionNames } = config;

  let rowFilters: MetricsViewFilterCond[] = [];
  if (isInitialTable && rowDimensionNames.length) {
    rowFilters = [
      {
        name: rowDimensionNames[0],
        in: rowDimensionAxes[rowDimensionNames[0]].slice(0, yLimit),
      },
    ];
  }

  const colFilters = colDimensionNames.map((colDimensionName) => {
    return {
      name: colDimensionName,
      in: colDimensionAxes?.[colDimensionName].slice(0, xLimit),
    };
  });

  const filters = {
    include: [...colFilters, ...rowFilters],
    exclude: [],
  };

  return filters;
}

/**
 * Create a nested accessor for a cell in the table.
 * This is used to map the cell data to the table data.
 *
 * Column names are converted to c0, c1, c2, etc.
 * Column values are converted to v0, v1, v2, etc.
 * Measure names are converted to m0, m1, m2, etc.
 */
export function getAccessorForCell(
  colDimensionNames: string[],
  colValuesIndexMaps: Map<string, number>[],
  numMeasures: number,
  cell: { [key: string]: string | number },
) {
  // TODO: Check for undefineds
  const nestedColumnValueAccessor = colDimensionNames
    .map((colName, i) => {
      let accessor = `c${i}`;

      const colValue = cell[colName] as string;
      const colValueIndex = colValuesIndexMaps[i].get(colValue);
      accessor += `v${colValueIndex}`;

      return accessor;
    })
    .join("_");

  return Array(numMeasures)
    .fill(null)
    .map((_, i) => `${nestedColumnValueAccessor}m${i}`);
}

/***
 * Create nested and grouped column definitions for pivot table
 */
function createColumnDefinitionForDimensions(
  dimensionNames: string[],
  headers: Record<string, string[]>,
  leafData: ColumnDef<PivotDataRow>[],
): ColumnDef<PivotDataRow>[] {
  const colValuesIndexMaps = dimensionNames?.map((colDimension) =>
    createIndexMap(headers[colDimension]),
  );

  const levels = dimensionNames.length;
  // Recursive function to create nested columns
  function createNestedColumns(
    level: number,
    colValuePair: { [key: string]: string },
  ): ColumnDef<PivotDataRow>[] {
    if (level === levels) {
      const accessors = getAccessorForCell(
        dimensionNames,
        colValuesIndexMaps,
        leafData.length,
        colValuePair,
      );

      // Base case: return leaf columns
      return leafData.map((leaf, i) => ({
        ...leaf,
        // Change accessor key to match the nested column structure
        accessorKey: accessors[i],
      }));
    }

    // Recursive case: create nested headers
    const headerValues = headers[dimensionNames?.[level]];
    return headerValues?.map((value) => ({
      header: value,
      columns: createNestedColumns(level + 1, {
        ...colValuePair,
        [dimensionNames[level]]: value,
      }),
    }));
  }

  // Start the recursion
  return createNestedColumns(0, {});
}

export function getColumnDefForPivot(
  config: PivotDataStoreConfig,
  columnDimensionAxes: Record<string, string[]> | undefined,
) {
  const IsNested = true;

  // TODO: Simplify function calls
  const measures = getMeasuresInPivotColumns(config.pivot, config.allMeasures);
  const rowDimensions = getDimensionsInPivotRow(
    config.pivot,
    config.allDimensions,
  );
  const colDimensions = getDimensionsInPivotColumns(
    config.pivot,
    config.allDimensions,
  );

  let rowDimensionsForColumnDef = rowDimensions;
  let nestedLabel: string;
  if (IsNested) {
    rowDimensionsForColumnDef = rowDimensions.slice(0, 1);
    nestedLabel = rowDimensions.map((d) => d.label || d.name).join(" > ");
  }
  const rowDefinitions: ColumnDef<PivotDataRow>[] =
    rowDimensionsForColumnDef.map((d) => {
      return {
        accessorKey: d.name,
        header: nestedLabel,
        cell: ({ row, getValue }) =>
          cellComponent(PivotExpandableCell, {
            value: getValue(),
            row,
          }),
      };
    });

  const leafColumns: ColumnDef<PivotDataRow>[] = measures.map((m) => {
    return {
      accessorKey: m.name as string,
      header: m.label || m.name,
      cell: (info) => info.getValue(),
    };
  });

  const groupedColDef = createColumnDefinitionForDimensions(
    (colDimensions.map((d) => d.column) as string[]) || [],
    columnDimensionAxes || {},
    leafColumns,
  );

  return [...rowDefinitions, ...groupedColDef];
}
