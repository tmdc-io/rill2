import type { PivotDataRow, PivotDataStoreConfig } from "./types";
import { createIndexMap, getAccessorForCell } from "./pivot-utils";
import type { V1MetricsViewAggregationResponseDataItem } from "@rilldata/web-common/runtime-client";

/**
 * Create a barebone table with row  headers.
 * This is used to render the table skeleton before cell data is fetched.
 */
export function createTableWithAxes(
  anchorDimensionName: string,
  rowDimensionValues: string[] | undefined,
) {
  let data: PivotDataRow[] = [];

  if (anchorDimensionName && rowDimensionValues && rowDimensionValues?.length) {
    data = rowDimensionValues?.map((value) => {
      return {
        [anchorDimensionName]: value,
      };
    });
  }

  return data;
}

/**
 * During the phase when queries are still being resolved, we don't have enough
 * information about the values present in an expanded group
 * For now fill it with empty values if there are more than one row dimensions
 */
export function prepareNestedPivotData(
  data: PivotDataRow[],
  dimensions: string[],
  i = 1,
) {
  if (dimensions.slice(i).length > 0) {
    data.forEach((row) => {
      if (!row.subRows) {
        row.subRows = [{ [dimensions[0]]: "LOADING_CELL" }];
      }

      prepareNestedPivotData(row.subRows, dimensions, i + 1);
    });
  }
}

/***
 *
 * The table expects pivoted data for a single row to be mapped to a single object
 * This functions reduces the API response from
 *
 * +-------------------------------+    +------------------------------------------------+
 * | {                             |    | {                                              |
 * |  "row_dimension: value,       |    |   "row_dimension": value,                      |
 * |  "column_dimension: value_1,  | => |   "column_dimension_value_1": measure_value_1, |
 * |  "measure": measure_value_1   |    |   "column_dimension_value_2": measure_value_2  |
 * | },                            |    | }                                              |
 * | {                             |    +------------------------------------------------+
 * |  "row_dimension: value,       |    Transformed to a Tanstack table readable format
 * |  "column_dimension: value_2,  |
 * |  "measure": measure_value_2   |
 * | }                             |
 * +-------------------------------+
 */
export function reduceTableCellDataIntoRows(
  config: PivotDataStoreConfig,
  anchorDimensionName: string,
  anchorDimensionRowValues: string[],
  columnDimensionAxes: Record<string, string[]>,
  tableData: PivotDataRow[],
  cellData: V1MetricsViewAggregationResponseDataItem[],
) {
  const colDimensionNames = config.colDimensionNames;

  // For simple tables, return the cell data as is
  if (!anchorDimensionName) {
    return cellData;
  }

  /**
   * Create a map of row dimension values to their index in the row dimension axes.
   * This way we can apply the sort order on the row dimension axes and sort the cell
   * data using the dimension values. For O(n) pass on the cells (n = number of cells),
   * we can reduce them optimally into row objects.
   */
  const rowDimensionIndexMap = createIndexMap(anchorDimensionRowValues);

  const colValuesIndexMaps = colDimensionNames.map((colDimensionName) =>
    createIndexMap(columnDimensionAxes[colDimensionName]),
  );

  cellData?.forEach((cell) => {
    const rowDimensionValue = cell[anchorDimensionName] as string;
    const rowIndex = rowDimensionIndexMap.get(rowDimensionValue);
    if (rowIndex === undefined) {
      console.warn(`Dimension value not found in row dimension axes`);
      return;
    }
    const row = tableData[rowIndex];

    const accessors = getAccessorForCell(
      colDimensionNames,
      colValuesIndexMaps,
      config.measureNames.length,
      cell,
    );

    if (row) {
      accessors.forEach((accessor, i) => {
        row[accessor] = cell[config.measureNames[i]] as string | number;
      });
    }
  });

  return tableData;
}
