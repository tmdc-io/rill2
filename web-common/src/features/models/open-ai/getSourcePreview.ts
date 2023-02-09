import {
  MutationFunction,
  useMutation,
  UseMutationOptions,
} from "@sveltestack/svelte-query";
import { RpcStatus, runtimeServiceQuery } from "../../../runtime-client";

export interface GetSourcePreviewRequest {
  instanceId: string;
  sourceName: string;
}

export interface GetSourcePreviewResponse {
  sourcePreview: string;
}

export const useGetSourcePreview = <
  TError = RpcStatus,
  TContext = unknown
>(options?: {
  mutation?: UseMutationOptions<
    Awaited<Promise<GetSourcePreviewResponse>>,
    TError,
    { data: GetSourcePreviewRequest },
    TContext
  >;
}) => {
  const { mutation: mutationOptions } = options ?? {};

  const mutationFn: MutationFunction<
    Awaited<Promise<GetSourcePreviewResponse>>,
    { data: GetSourcePreviewRequest }
  > = async (props) => {
    const { data } = props ?? {};

    const schema = await getSourceSchema(data.instanceId, data.sourceName);
    const exampleDataHead = await getSourceExampleDataHead(
      data.instanceId,
      data.sourceName
    );
    const exampleDataTail = await getSourceExampleDataTail(
      data.instanceId,
      data.sourceName
    );

    return {
      sourcePreview:
        schema + "\n#\n" + exampleDataHead + "\n#\n" + exampleDataTail,
    };
  };

  return useMutation<
    Awaited<Promise<GetSourcePreviewResponse>>,
    TError,
    { data: GetSourcePreviewRequest },
    TContext
  >(mutationFn, mutationOptions);
};

async function getAllSourceNames(instanceId: string): Promise<string[]> {
  const resp = await runtimeServiceQuery(instanceId, {
    sql: "select distinct table_name from information_schema.columns where table_schema = 'main';",
    priority: 1,
  });

  return resp.data.map((row) => row["table_name"]);
}

async function getSourceSchema(
  instanceId: string,
  sourceName: string
): Promise<string> {
  const header = `#\tSchema for table \`${sourceName}\`\n#\n`;

  const resp = await runtimeServiceQuery(instanceId, {
    sql: `select table_name, column_name, data_type from information_schema.columns where table_schema = 'main' and table_name = '${sourceName}';`,
    priority: 1,
  });

  // input: {'table_name': TABLE_NAME, 'column_name': COLUMN_NAME, 'data_type': DATA_TYPE}[]
  // output: # Table_Name(column_name_1: DATA_TYPE, column_name_2: DATA_TYPE)\n # Table_Name(...)\n
  const schema = resp.data
    .map((row) => {
      return `#\t\t${row["column_name"]}: ${row["data_type"]}`;
    })
    .join("\n");

  return header + schema;
}

async function getSourceExampleDataHead(
  instanceId: string,
  sourceName: string
): Promise<string> {
  const n_rows = 5;
  const resp = await runtimeServiceQuery(instanceId, {
    sql: `select rowid, * from ${sourceName} order by rowid limit 5;`,
    priority: 1,
  });
  const exampleData = formatResponse(resp);
  const header = `#\tFirst ${n_rows} rows of table \`${sourceName}\`\n#\n`;
  return header + exampleData;
}

async function getSourceExampleDataTail(
  instanceId: string,
  sourceName: string
): Promise<string> {
  const n_rows = 5;
  const resp = await runtimeServiceQuery(instanceId, {
    sql: `select rowid, * from ${sourceName} order by rowid desc limit ${n_rows};`,
    priority: 1,
  });
  const exampleData = formatResponse(resp);
  const header = `#\tLast ${n_rows} rows of table \`${sourceName}\`\n#\n`;
  return header + exampleData;
}

function formatResponse(resp: any): string {
  const colNames = resp.meta.fields.map((field) => field.name);

  // Find max width for each column.
  const colWidths = resp.data.reduce((acc, row) => {
    colNames.forEach((colName, i) => {
      const colWidth = acc[i] ?? 0;
      const item = row[colName] ?? "";
      const itemWidth = item.toString().length;
      acc[i] = Math.max(colWidth, itemWidth, colNames[i].length);
    });
    return acc;
  }, [] as number[]);

  // Format the column names.
  const formattedColNames =
    "#\t" +
    colNames
      .map((colName, i) => colName.padStart(colWidths[i], " "))
      .join("\t");

  // Format each row of the table.
  const formattedData = resp.data
    .map((row) => {
      return (
        "#\t" +
        colNames
          .map((colName, j) => {
            const item = row[colName] ?? "";
            return item.toString().padStart(colWidths[j], " ");
          })
          .join("\t")
      );
    })
    .join("\n");

  return formattedColNames + "\n" + formattedData;
}
