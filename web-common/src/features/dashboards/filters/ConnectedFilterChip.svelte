<!-- @component
The main feature-set component for dashboard filters
 -->
<script lang="ts">
  import { RemovableListChip } from "@rilldata/web-common/components/chip";
  import { defaultChipColors } from "@rilldata/web-common/components/chip/chip-types";
  // import Filter from "@rilldata/web-common/components/icons/Filter.svelte";
  // import FilterRemove from "@rilldata/web-common/components/icons/FilterRemove.svelte";
  import { cancelDashboardQueries } from "@rilldata/web-common/features/dashboards/dashboard-queries";
  import {
    useMetaQuery,
    useModelHasTimeSeries,
  } from "@rilldata/web-common/features/dashboards/selectors";
  import type {
    MetricsViewDimension,
    MetricsViewFilterCond,
    V1MetricsViewFilter,
  } from "@rilldata/web-common/runtime-client";
  import { createQueryServiceMetricsViewToplist } from "@rilldata/web-common/runtime-client";
  import { getMapFromArray } from "@rilldata/web-local/lib/util/arrayUtils";
  import { useQueryClient } from "@tanstack/svelte-query";
  // import { flip } from "svelte/animate";
  // import { fly } from "svelte/transition";
  import { runtime } from "../../../runtime-client/runtime-store";
  import {
    MetricsExplorerEntity,
    metricsExplorerStore,
  } from "../dashboard-stores";
  import { getDisplayName } from "./getDisplayName";
  import { afterUpdate, beforeUpdate, onDestroy, onMount } from "svelte";

  export let metricViewName;
  export let dimensionName;

  const queryClient = useQueryClient();

  let metricsExplorer: MetricsExplorerEntity;
  $: metricsExplorer = $metricsExplorerStore.entities[metricViewName];

  let includeValues: Array<MetricsViewFilterCond>;
  $: includeValues = metricsExplorer?.filters.include;

  let excludeValues: Array<MetricsViewFilterCond>;
  $: excludeValues = metricsExplorer?.filters.exclude;

  $: metaQuery = useMetaQuery($runtime.instanceId, metricViewName);
  let dimensions: Array<MetricsViewDimension>;
  $: dimensions = $metaQuery.data?.dimensions;

  function clearFilterForDimension(dimensionId, include: boolean) {
    cancelDashboardQueries(queryClient, metricViewName);
    metricsExplorerStore.clearFilterForDimension(
      metricViewName,
      dimensionId,
      include
    );
  }

  // function isFiltered(filters: V1MetricsViewFilter): boolean {
  //   if (!filters) return false;
  //   return filters.include.length > 0 || filters.exclude.length > 0;
  // }

  let topListQuery;
  let searchText = "";
  let searchedValues = [];
  // let activeDimensionName;

  $: metricTimeSeries = useModelHasTimeSeries(
    $runtime.instanceId,
    metricViewName
  );
  $: hasTimeSeries = $metricTimeSeries.data;

  $: addNull = "null".includes(searchText);

  $: if (searchText == "") {
    searchedValues = [];
  } else {
    let topListParams = {
      dimensionName,
      limit: "15",
      offset: "0",
      sort: [],
      filter: {
        include: [
          {
            name: dimensionName,
            in: addNull ? [null] : [],
            like: [`%${searchText}%`],
          },
        ],
        exclude: [],
      },
    };

    if (hasTimeSeries) {
      topListParams = {
        ...topListParams,
        ...{
          timeStart: metricsExplorer?.selectedTimeRange?.start,
          timeEnd: metricsExplorer?.selectedTimeRange?.end,
        },
      };
    }
    // Use topList API to fetch the dimension names
    // We prune the measure values and use the dimension labels for the filter
    topListQuery = createQueryServiceMetricsViewToplist(
      $runtime.instanceId,
      metricViewName,
      topListParams
    );
  }

  $: if (!$topListQuery?.isFetching && searchText != "") {
    const topListData = $topListQuery?.data?.data ?? [];
    searchedValues = topListData.map((datum) => datum[dimensionName]) ?? [];
  }

  // $: hasFilters = isFiltered(metricsExplorer?.filters);

  /** prune the values and prepare for templating */
  // let currentDimensionFilters = [];
  // $: if (includeValues && excludeValues && dimensions) {
  //   const dimensionIdMap = getMapFromArray(
  //     dimensions,
  //     (dimension) => dimension.name
  //   );
  //   const currentDimensionIncludeFilters = includeValues.map(
  //     (dimensionValues) => ({
  //       name: dimensionValues.name,
  //       label: getDisplayName(dimensionIdMap.get(dimensionValues.name)),
  //       selectedValues: dimensionValues.in,
  //       filterType: "include",
  //     })
  //   );
  //   const currentDimensionExcludeFilters = excludeValues.map(
  //     (dimensionValues) => ({
  //       name: dimensionValues.name,
  //       label: getDisplayName(dimensionIdMap.get(dimensionValues.name)),
  //       selectedValues: dimensionValues.in,
  //       filterType: "exclude",
  //     })
  //   );
  //   currentDimensionFilters = [
  //     ...currentDimensionIncludeFilters,
  //     ...currentDimensionExcludeFilters,
  //   ];
  //   // sort based on name to make sure toggling include/exclude is not jarring
  //   currentDimensionFilters.sort((a, b) => (a.name > b.name ? 1 : -1));
  //   console.log({ currentDimensionFilters });
  // }

  $: console.log("metricViewName", metricViewName);
  $: console.log("dimensionName", dimensionName);
  $: console.log("dimensions", dimensions);

  $: dimensionIdMap = getMapFromArray(
    dimensions ?? [],
    (dimension) => dimension.name
  );

  $: label = getDisplayName(dimensionIdMap.get(dimensionName));

  // if this dimensionName is included among the `includeValues`,
  // then this dimension is in include mode.
  $: isInclude =
    (includeValues?.filter(
      (dimensionValues) => dimensionValues.name === dimensionName
    )?.length ?? 0) > 0;

  // NOTE: there should only ever be includeValues or excludeValues
  // additionally, filtering by `dimensionName` should only return
  // a list with one MetricsViewFilterCond, so we can grab the 0th
  // entry in the list and access its `in` prop.
  // Additionally, we only want to update `selectedValues` when the
  // menu is NOT active. This prevents a re-render when the
  // include/exclude status of an individual item is toggled.
  // let selectedValues = [];

  let selectedValues: string[] = [];

  $: if (!active) {
    selectedValues =
      (isInclude ? includeValues : excludeValues)?.filter(
        (dimensionValues) => dimensionValues.name === dimensionName
      )[0]?.in ?? [];
    console.log("updating selectedValues", dimensionName);
  }

  // let selectedValuesSet: Set<string> = new Set();

  // const unionSelectedValues = (newValues: string[]) => {
  //   selectedValuesSet = new Set([...selectedValuesSet, ...newValues]);
  //   selectedValues = [...selectedValuesSet];
  //   // sort based on name to make sure toggling include/exclude is not jarring
  //   selectedValues.sort((a, b) => (a > b ? 1 : -1));
  // };

  // $: unionSelectedValues(
  //   (isInclude ? includeValues : excludeValues)?.filter(
  //     (dimensionValues) => dimensionValues.name === dimensionName
  //   )[0]?.in ?? []
  // );

  // let selectedValues =
  //   (isInclude ? includeValues : excludeValues)?.filter(
  //     (dimensionValues) => dimensionValues.name === dimensionName
  //   )[0]?.in ?? [];

  // $: if (!active) {
  //   console.log("active changed SHOULD BE FALSE", dimensionName, active);
  // }
  // $: console.log("CHANGED active", active, " for dim ", dimensionName);

  // $: console.log("CHANGED dimensionName", dimensionName);
  // $: console.log("CHANGED active", active, " for dim ", dimensionName);
  // $: console.log(
  //   "CHANGED selectedValues",
  //   selectedValues.toString().slice(0, 0),
  //   " for dim ",
  //   dimensionName
  // );

  // $: console.log(
  //   "CHANGED includeValues",
  //   includeValues,
  //   " for dim ",
  //   dimensionName
  // );
  // $: console.log(
  //   "CHANGED excludeValues",
  //   excludeValues,
  //   " for dim ",
  //   dimensionName
  // );

  // beforeUpdate(() =>
  //   console.log(
  //     "ConnectecFilter updated: ",
  //     dimensionName
  //     // selectedValues
  //   )
  // );
  // afterUpdate(() =>
  //   console.log(
  //     "ConnectecFilter updated: ",
  //     dimensionName
  //     // selectedValues
  //   )
  // );
  // onMount(() => console.log("ConnectecFilter MOUNT: ", dimensionName));
  // onDestroy(() => console.log("ConnectecFilter DESTROY: ", dimensionName));

  function toggleDimensionValue(event) {
    cancelDashboardQueries(queryClient, metricViewName);
    metricsExplorerStore.toggleFilter(
      metricViewName,
      dimensionName,
      event.detail
    );
  }

  function toggleFilterMode() {
    cancelDashboardQueries(queryClient, metricViewName);
    metricsExplorerStore.toggleFilterMode(metricViewName, dimensionName);
  }

  const excludeChipColors = {
    bgBaseClass: "bg-gray-100 dark:bg-gray-700",
    bgHoverClass: "bg-gray-200 dark:bg-gray-600",
    textClass: "ui-copy",
    bgActiveClass: "bg-gray-200 dark:bg-gray-600",
    outlineClass: "outline-gray-400 dark:outline-gray-500",
  };

  let active: boolean;

  $: console.log("Removable list chip active", active);
</script>

<RemovableListChip
  bind:active
  on:toggle={() => toggleFilterMode()}
  on:remove={() =>
    clearFilterForDimension(metricViewName, isInclude ? true : false)}
  on:apply={(event) => toggleDimensionValue(event)}
  on:search={(event) => {
    searchText = event.detail;
  }}
  typeLabel="dimension"
  name={isInclude ? label : `Exclude ${label}`}
  excludeMode={isInclude ? false : true}
  colors={isInclude ? defaultChipColors : excludeChipColors}
  {selectedValues}
  {searchedValues}
>
  <svelte:fragment slot="body-tooltip-content">
    Click to edit the the filters in this dimension
  </svelte:fragment>
</RemovableListChip>
