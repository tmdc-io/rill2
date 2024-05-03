from utils import * 

replace_chars = lambda s: s.replace('`', '').replace('TABLE.', '').replace('{', ''). \
    replace('}', '').replace('$', '')


def create_board_yaml(lens_meta):
    lens_name = lens_meta['name']
    dir_name = lens_name
    transformed_lens_meta = measure_dimension_sqls(lens_meta)
    for v in lens_meta['tables']:
        t_v_name = v['name']
        meta_info = v.get('meta', {})
        if 'export_to_board' in meta_info.keys() and meta_info['export_to_board'] is True:
            print(f"Generating board yaml for `{v['type']}` - {t_v_name}")
            if v['type'] == 'table' and v['public']:
                load_data_obj = {"dimensions": [], "name": t_v_name}
                dimensions = []
                measures = []
                additional_kv = None
                for d in v.get('dimensions', []):
                    d_name = d['name'].split('.')[1]
                    dimensions.append(
                        {"label": d['shortTitle'],
                         "description": d.get('description', d['title']),
                         "name": d_name,
                         "property": d_name
                         }
                    )
                    load_data_obj["dimensions"].append((d_name, d['type']))
                for m in v.get("measures", []):
                    m_name = m['name'].split('.')[1]
                    m_sql = replace_chars(m['sql'])
                    measures.append(
                        {
                            "label": m['shortTitle'],
                            "name": m_name,
                            "description": m.get('description', m['title']),
                            "expression": get_sql_expression(agg_type=m['aggType'], sql=m_sql, table_name=v['name'],
                                                             is_prefix=False)
                        }
                    )

                if 'board' in meta_info.keys():
                    additional_kv = meta_info['board']

                # create board yaml
                dump_board_yaml(lens_name=dir_name, dimensions=dimensions, measures=measures,
                                additional_kv=additional_kv, data_obj=load_data_obj)

            elif meta_info['export_to_board'] and v['type'] == 'view':
                load_data_obj = {"dimensions": [], "name": t_v_name}
                prefix = False
                dimensions = []
                measures = []
                additional_kv = None
                existing_dimensions = []
                if v.get('dimensions', []) != 0:
                    for d in v['dimensions']:
                        d_name = d['name'].split('.')[1]
                        d_alias = d['aliasMember']
                        prefix = True if d_alias.split('.')[0] in d_name else False
                        dimensions.append({
                            "label": d['shortTitle'],
                            "description": d.get('description', d['title']),
                            "name": d_name,
                            "property": d_name
                        })
                        existing_dimensions.append(d_name)
                        load_data_obj["dimensions"].append((d_name, d['type']))
                    for m in v.get('measures', []):
                        m_name = m['name'].split('.')[1]
                        m_alias = m['aliasMember']
                        m_sql = replace_chars(transformed_lens_meta[m_alias]['sql'])
                        m_query = f"SELECT {m_sql}"
                        m_d_columns = Parser(m_query).columns  # measure sql columns
                        is_skip = False  # Assume we don't skip by default
                        is_prefix = True if m_alias.split('.')[0] in m_name else False
                        prefix = True if m_alias.split('.')[0] in m_name else False
                        for m_dim_name in m_d_columns:
                            if is_prefix:
                                m_d_name_with_table = f"{m_alias.split('.')[0]}_{m_dim_name}"
                            else:
                                m_d_name_with_table = m_dim_name
                            if m_d_name_with_table not in existing_dimensions:
                                is_skip = True  # If any dimension is missing, we set to skip
                                print(
                                    f"Skipping measure - `{m_name}`, its dependent dimension - `{m_dim_name}` is missing "
                                    f"in view - `{v['name']}`")
                                break
                        if is_skip is False:
                            measures.append({
                                "label": m['shortTitle'],
                                "name": m_name,
                                "description": m.get('description', m['title']),
                                "expression": get_sql_expression(agg_type=m['aggType'], sql=m_sql,
                                                                 table_name=m_alias.split('.')[0],
                                                                 is_prefix=is_prefix)
                            })
                    if 'board' in meta_info.keys():
                        additional_kv = meta_info['board']
                        if prefix:
                            additional_kv['timeseries'] = additional_kv['timeseries'].replace(".", "_")
                        else:
                            additional_kv['timeseries'] = additional_kv['timeseries'].split(".")[1]

                    dump_board_yaml(lens_name=dir_name, dimensions=dimensions, measures=measures,
                                    additional_kv=additional_kv, data_obj=load_data_obj)
                else:
                    print(f"No dimension found  for view - `{v['name']}`")
        else:
            print(f"Skipping `{v['type']}` `{t_v_name}` as it's meta does not contains `export_to_board` key or set "
                  f"to `False`")
    return

meta = get_lens_meta()
create_board_yaml(meta)