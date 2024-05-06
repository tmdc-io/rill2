import json
import os

import requests
import yaml
from sql_metadata import Parser

from constants import LENS2_BASE_URL, LENS2_NAME, DATAOS_RUN_AS_APIKEY


def get_env_or_throw(env_name):
    """Get the value of an environment variable or raise an error."""
    value = os.getenv(env_name)
    if value is None:
        raise ValueError(f"The environment variable '{env_name}' is not set.")
    return value


def make_api_call_with_headers(api_url, headers):
    """Make an API call with headers."""
    response = requests.get(api_url, headers=headers)
    return response


def get_lens_meta():
    lens2_base_url = get_env_or_throw(LENS2_BASE_URL)
    lens2_name = get_env_or_throw(LENS2_NAME)
    apikey = get_env_or_throw(DATAOS_RUN_AS_APIKEY)

    url = f"{lens2_base_url}/{lens2_name}/v2/meta?showPrivate=true"
    headers = {"apikey": apikey}
    return make_api_call_with_headers(url, headers=headers).json()


def measure_dimension_sqls(lens_meta):
    transformed_meta = {}
    for t in lens_meta['tables']:
        if t['type'] == 'table':
            for m in t['measures']:
                transformed_meta[m['name']] = {'sql': m['sql'],
                                               'info': m}
            for d in t['dimensions']:
                transformed_meta[d['name']] = {'sql': d['sql'],
                                               'info': d}
    return transformed_meta


def replace_with_dict(string, replacement_dict):
    for key, value in replacement_dict.items():
        string = string.replace(key, value)
    return string


def get_sql_expression(agg_type=None, sql=None, table_name=None, is_prefix=None):
    # SQL Manipulation
    cols = Parser(f"SELECT {sql}").columns
    cols = {col: f"{table_name}_{col}" for col in cols}
    if is_prefix:
        sql = replace_with_dict(sql, cols)

    if agg_type in ["countDistinct", "countDistinctApprox"]:
        return f"COUNT(DISTINCT {sql})"
    elif agg_type in ["count", "sum", "avg", "min", "max"]:
        return f"{agg_type}({sql})"
    elif agg_type in ["string", "time", "boolean", "number"]:
        return f"{sql}"


def dump_board_yaml(lens_name=None, measures=None, dimensions=None, additional_kv=None,
                    data_obj=None):  # assumption data_dir = /etc/dataos/work/data/tableOrViewName

    # Models
    model_path = os.path.join(os.getcwd(), "models")
    if not os.path.exists(model_path):
        os.makedirs(model_path, exist_ok=False)

    # Sources
    source_path = os.path.join(os.getcwd(), "sources")
    if not os.path.exists(source_path):
        os.makedirs(source_path, exist_ok=False)

    # Dashboard
    dashboard_path = os.path.join(os.getcwd(), f"dashboards")
    if not os.path.exists(dashboard_path):
        os.makedirs(dashboard_path, exist_ok=False)


    # Required rill file
    rill_path = os.path.join(os.getcwd(), "rill.yaml")
    with open(rill_path, 'w') as file:
        board_data = {"compiler": "rill",
                      "name": lens_name}
        file.write(yaml.dump(board_data, default_flow_style=False, sort_keys=False))

    # source files
    t_v_path = os.path.join(source_path, f"{data_obj['name']}_source.yaml")
    with open(t_v_path, 'w') as file:
        source = {
            "connector": "local_file",
            "path": f"/etc/dataos/work/data/{data_obj['name']}/*.parquet",
            "lens": {
                "baseUri": get_env_or_throw(LENS2_BASE_URL),
                "name": lens_name,
                "query": {
                    "dimensions": [f"{data_obj['name']}.{dim[0]}" for dim in data_obj['dimensions']],
                    "batch": 50000,
                    "start": 0,
                    "end": -1
                },
                "apikey": get_env_or_throw(DATAOS_RUN_AS_APIKEY)
            }
        }
        file.write(yaml.dump(source, default_flow_style=False, sort_keys=False))

    # models yaml
    with open(os.path.join(model_path, f"{data_obj['name']}_model.sql"), 'w') as file:
        file.write(f"""SELECT * FROM {data_obj['name']}_source""")

    # dashboards yaml
    dashboard_data = {
        "title": data_obj['name'],
        "model": f"{data_obj['name']}_model",
    }
    if additional_kv:
        dashboard_data.update(additional_kv)
    dashboard_data.update(
        {
            "dimensions": dimensions,
            "measures": measures
        }
    )
    with open(os.path.join(dashboard_path, f"{data_obj['name']}.yaml"), 'w') as file:
        file.write(yaml.dump(dashboard_data, default_flow_style=False, sort_keys=False))
