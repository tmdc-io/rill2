<script lang="ts">
    import Select from "@rilldata/web-common/components/forms/Select.svelte"
    import Input from "@rilldata/web-common/components/forms/Input.svelte"
    import { getStateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";
    import { 
        createQueryServiceMetricsViewAggregation
    } from "@rilldata/web-common/runtime-client";

    $: dimension = ''
    $: measure = ''
    $: havingClause = ''

    const {
        selectors: {
            dimensions: { allDimensions },
            measures: { allMeasures },
            timeRangeSelectors: {
                timeControlsState
            }
        },
        metricsViewName,
        runtime,
        queryClient,
        dashboardStore
    } = getStateManagers();

    let dimensionQuery: any[] = [];
    $: if (dimension.length) {
        dimensionQuery = []
        dimensionQuery.push({ name: dimension });
    }
    let havingQuery: any[] = [];
    $: if (havingClause.length) {
        havingQuery = [];
        havingQuery.push($allMeasures.filter(m => m.name === measure)[0].expression + " " + havingClause);
    }

    $: alertData = createQueryServiceMetricsViewAggregation(
        $runtime.instanceId,
        $metricsViewName,
        {
            dimensions: dimensionQuery,
            measures: [{ name: measure }],
            filter: {
                exclude: [],
                include: $dashboardStore.filters.include,
                having: havingQuery
            },
            timeStart: $timeControlsState.timeStart,
            timeEnd: $timeControlsState.timeEnd,
            limit: "20",
            offset: "0",
        },
        {
            query: {
                enabled: true,
                queryClient,
            },
        }
    )
    let headers: any;
    $: headers = $alertData.data?.schema?.fields?.map(d => d.name)
    let rows: any;
    $: rows = $alertData.data?.data?.map(d => Object.values(d))
    $: console.log(rows)

</script>
<div class="alertbase">
    <h2>Alert Prototype</h2>

    <Select
        bind:value={measure}
        id="measure"
        label="Measure"
        options={$allMeasures.map((m) => ({
            value: m.name,
        }))}
    />
    <Select
        bind:value={dimension}
        id="dimension"
        label="Dimension (Optional)"
        options={$allDimensions.map((dim) => ({
            value: dim.name,
        }))}
    />


    <Input id="having" bind:value={havingClause} label="Criteria" />
    <div class="table-container">
        {#if $alertData.isSuccess}
        <table>
            <thead>
                <tr>
                    {#each headers as head}
                        <th>{head}</th>
                    {/each}
                </tr>
            </thead>
            <tbody>
                {#each rows as row}
                <tr>
                    {#each row as value}
                        <td>{value}</td>
                    {/each}
                </tr>
                {/each}
            </tbody>
        </table>
        {/if}
    </div>

</div>

<style>
    h2 {
        margin-bottom: 20px;
    }
    table {
        display: table;
        width: 100%;
        border-collapse: separate;
        border-spacing: 0px;
        border-left: 1px solid rgb(217, 217, 217);
        border-top: 1px solid rgb(217, 217, 217);
        max-height: 500px;
    }
    th {
        font-size: 14px;
        line-height: 1.71429rem;
        display: table-cell;
        vertical-align: inherit;
        text-align: left;
        color: rgb(64, 64, 64);
        position: sticky;
        top: 0px;
        z-index: 2;
        padding: 0px 8px 0px 16px;
        border-right: 1px solid rgb(217, 217, 217);
        height: 41px;
        font-weight: 600;
        background-color: rgb(250, 250, 250);
        border-bottom: 1px solid rgb(217, 217, 217);
    }
    tbody {
        display: table-row-group;
    }
    td {
        font-size: 14px;
        line-height: 20px;
        font-weight: 400;
        display: table-cell;
        vertical-align: inherit;
        text-align: left;
        color: rgb(64, 64, 64);
        padding: 0px 8px 0px 16px;
        border-bottom: 1px solid rgb(217, 217, 217);
        border-right: 1px solid rgb(217, 217, 217);
        height: 39px;
    }
    .alertbase {
        min-width: 450px;
    }
    .table-container {
        max-height: 500px;
        overflow: scroll;
        margin-top: 20px;
    }
</style>
