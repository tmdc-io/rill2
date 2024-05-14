from utils import *
import sys
import subprocess


def create_board_yaml(lens_meta):
    lens_name = lens_meta['name']
    dir_name = lens_name
    transformed_lens_meta = measure_dimension_sqls(lens_meta)
    start_iris_board = False
    for table_info in lens_meta['tables']:
        name = table_info['name']
        meta_info = table_info.get('meta', {})

        # check if board or iris exist in table meta
        if ('export_to_board' in meta_info.keys() and meta_info['export_to_board'] is True) \
                or ('export_to_iris' in meta_info.keys() and meta_info['export_to_iris'] is True):
            print(f"✅ Processing {table_info['type']}:{name}")

            if table_info['type'] == 'table' and table_info['public']:
                data_obj = {"dimensions": [], "name": name, "title": table_info['title']}
                dimensions = []
                measures = []
                for dimension in table_info.get('dimensions', []):
                    dim_name = dimension['name'].split('.')[1]
                    dimensions.append(
                        {"label": dimension['shortTitle'],
                         "description": dimension.get('description', dimension['title']),
                         "name": dim_name,
                         "property": dim_name
                         }
                    )
                    data_obj["dimensions"].append((dim_name, dimension['type']))
                for measure in table_info.get("measures", []):
                    measure_name = measure['name'].split('.')[1]
                    m_sql = replace_chars(measure['sql'])
                    measures.append(
                        {
                            "label": measure['shortTitle'],
                            "name": measure_name,
                            "description": measure.get('description', measure['title']),
                            "expression": get_sql_expression(agg_type=measure['aggType'], sql=m_sql, table_name=name,
                                                             is_prefix=False)
                        }
                    )

                additional_kv = {}
                excludes_dimension = []
                if 'board' in meta_info.keys() or 'iris' in meta_info.keys():
                    iris_board_meta = meta_info.get('board', {})
                    iris_board_meta.update(meta_info.get('iris', {}))
                    if "excludes" in iris_board_meta.keys():
                        excludes_dimension = [d.split(".")[1] for d in iris_board_meta['excludes']]
                        del iris_board_meta['excludes']
                    additional_kv['available_time_zones'] = lens_meta.get("timeZones", "")
                    additional_kv.update(iris_board_meta)
                    if 'timeseries' in additional_kv.keys():
                        additional_kv['timeseries'] = additional_kv.get('timeseries', '').split(".")[1]
                # create board yaml
                dump_board_yaml(lens_name=dir_name, dimensions=dimensions, measures=measures,
                                additional_kv=additional_kv, data_obj=data_obj, excludes_dimension=excludes_dimension)
                start_iris_board = True

            elif table_info['type'] == 'view' and table_info['public']:
                data_obj = {"dimensions": [], "name": name, "title": table_info['title']}
                is_prefix = False
                dimensions = []
                measures = []
                existing_dimensions = []
                if len(table_info.get('dimensions', [])) != 0:
                    for dimension in table_info['dimensions']:
                        dimension_name = dimension['name'].split('.')[1]
                        dimension_alias = dimension['aliasMember']
                        dimensions.append({
                            "label": dimension['shortTitle'],
                            "description": dimension.get('description', dimension['title']),
                            "name": dimension_name,
                            "property": dimension_name
                        })
                        existing_dimensions.append(dimension_name)
                        data_obj["dimensions"].append((dimension_name, dimension['type']))
                    for measure in table_info.get('measures', []):
                        measure_name = measure['name'].split('.')[1]
                        measure_alias = measure['aliasMember']
                        measure_sql = replace_chars(transformed_lens_meta[measure_alias]['sql'])
                        measure_query = f"SELECT {measure_sql}"
                        measure_dims_columns = Parser(measure_query).columns  # measure sql columns
                        is_skip = False  # Assume we don't skip by default
                        is_prefix = True if measure_alias.split('.')[0] in measure_name else False
                        for measure_dim_name in measure_dims_columns:
                            if is_prefix:
                                m_d_name_with_table = f"{measure_alias.split('.')[0]}_{measure_dim_name}"
                            else:
                                m_d_name_with_table = measure_dim_name
                            if m_d_name_with_table not in existing_dimensions:
                                is_skip = True  # If any dimension is missing, we set to skip
                                print(
                                    f"⚠️ Skipped measure - `{measure_name}`: `{table_info['name']}`")
                                break
                        if is_skip is False:
                            measures.append({
                                "label": measure['shortTitle'],
                                "name": measure_name,
                                "description": measure.get('description', measure['title']),
                                "expression": get_sql_expression(agg_type=measure['aggType'], sql=measure_sql,
                                                                 table_name=measure_alias.split('.')[0],
                                                                 is_prefix=is_prefix)
                            })

                    additional_kv = {}
                    excludes_dimension = []
                    if 'board' in meta_info.keys() or 'iris' in meta_info.keys():
                        iris_board_meta = meta_info.get('board', {})
                        iris_board_meta.update(meta_info.get('iris', {}))
                        if "excludes" in iris_board_meta.keys():
                            if is_prefix:
                                excludes_dimension = [d.replace(".", "_") for d in iris_board_meta['excludes']]
                            else:
                                excludes_dimension = [d.split(".")[1] for d in iris_board_meta['excludes']]
                            del iris_board_meta['excludes']
                        additional_kv['available_time_zones'] = lens_meta.get("timeZones", "")
                        additional_kv.update(iris_board_meta)
                        if 'timeseries' in additional_kv.keys():
                            if is_prefix:
                                additional_kv['timeseries'] = additional_kv.get('timeseries', '').replace(".", "_")
                            else:
                                additional_kv['timeseries'] = additional_kv.get('timeseries', '').split(".")[1]

                    dump_board_yaml(lens_name=dir_name, dimensions=dimensions, measures=measures,
                                    additional_kv=additional_kv, data_obj=data_obj,
                                    excludes_dimension=excludes_dimension)
                    start_iris_board = True
                else:
                    print(f"No dimension found  for view - `{table_info['name']}`")
        else:
            print(f"❌ Skipped {table_info['type']}:{name}")
    return start_iris_board

if __name__ == "__main__":
    meta = get_lens_meta()
    rill_start = create_board_yaml(meta)

    # delegate downstream
    if rill_start:
        subprocess.call(sys.argv[1:])