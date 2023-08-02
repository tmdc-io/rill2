<script lang="ts">
  import Body from "@rilldata/web-common/components/data-graphic/elements/Body.svelte";
  import SimpleDataGraphic from "@rilldata/web-common/components/data-graphic/elements/SimpleDataGraphic.svelte";
  import WithBisector from "@rilldata/web-common/components/data-graphic/functional-components/WithBisector.svelte";
  import WithRoundToTimegrain from "@rilldata/web-common/components/data-graphic/functional-components/WithRoundToTimegrain.svelte";
  import {
    Axis,
    Grid,
  } from "@rilldata/web-common/components/data-graphic/guides";
  import { ChunkedLine } from "@rilldata/web-common/components/data-graphic/marks";
  import { NumberKind } from "@rilldata/web-common/lib/number-formatting/humanizer-types";
  import { previousValueStore } from "@rilldata/web-common/lib/store-utils";
  import type { V1TimeGrain } from "@rilldata/web-common/runtime-client";
  import { extent } from "d3-array";
  import { cubicOut } from "svelte/easing";
  import { writable } from "svelte/store";
  import { fly } from "svelte/transition";
  import MeasureValueMouseover from "./MeasureValueMouseover.svelte";
  import {
    getBisectedTimeFromCordinates,
    getOrderedStartEnd,
    niceMeasureExtents,
  } from "./utils";
  import {
    TimeRangePreset,
    TimeRoundingStrategy,
  } from "../../../lib/time/types";
  import { getContext } from "svelte";
  import { contexts } from "@rilldata/web-common/components/data-graphic/constants";
  import type { ScaleStore } from "@rilldata/web-common/components/data-graphic/state/types";
  import { TIME_GRAIN } from "@rilldata/web-common/lib/time/config";
  import MeasureScrub from "@rilldata/web-common/features/dashboards/time-series/MeasureScrub.svelte";
  import {
    metricsExplorerStore,
    useDashboardStore,
  } from "@rilldata/web-common/features/dashboards/dashboard-stores";

  export let metricViewName: string;
  export let width: number = undefined;
  export let height: number = undefined;
  export let xMin: Date = undefined;
  export let xMax: Date = undefined;
  export let yMin: number = undefined;
  export let yMax: number = undefined;

  export let timeGrain: V1TimeGrain;

  export let showComparison = false;
  export let data;
  export let xAccessor = "ts";
  export let labelAccessor = "label";
  export let yAccessor = "value";
  export let mouseoverValue;
  export let hovered = false;

  // control point for scrub functionality.
  export let scrubbing = false;
  export let scrubStart;
  export let scrubEnd;

  export let mouseoverFormat: (d: number | Date | string) => string = (v) =>
    v.toString();
  export let mouseoverTimeFormat: (d: number | Date | string) => string = (v) =>
    v.toString();
  export let numberKind: NumberKind = NumberKind.ANY;
  export let tweenProps = { duration: 400, easing: cubicOut };

  $: dashboardStore = useDashboardStore(metricViewName);

  const xScale = getContext(contexts.scale("x")) as ScaleStore;

  let justCreatedScrub = false;

  // $: console.log(scrubEnd, yAccessor);

  $: hasSubrangeSelected = Boolean(scrubStart && scrubEnd);

  $: mainLineColor =
    scrubbing || hasSubrangeSelected
      ? "hsla(217, 10%, 60%, 1)"
      : "hsla(217,60%, 55%, 1)";

  $: areaColor =
    scrubbing || hasSubrangeSelected
      ? "hsla(225, 20%, 80%, .2)"
      : "hsla(217,70%, 80%, .4)";

  $: [xExtentMin, xExtentMax] = extent(data, (d) => d[xAccessor]);
  $: [yExtentMin, yExtentMax] = extent(data, (d) => d[yAccessor]);
  let comparisonExtents;

  /** if we are making a comparison, factor this into the extents calculation.*/
  $: if (showComparison) {
    comparisonExtents = extent(data, (d) => d[`comparison.${yAccessor}`]);

    yExtentMin = Math.min(yExtentMin, comparisonExtents[0] || yExtentMin);
    yExtentMax = Math.max(yExtentMax, comparisonExtents[1] || yExtentMax);
  }

  $: [internalYMin, internalYMax] = niceMeasureExtents(
    [
      yMin !== undefined ? yMin : yExtentMin,
      yMax !== undefined ? yMax : yExtentMax,
    ],
    6 / 5
  );

  $: internalXMin = xMin || xExtentMin;
  $: internalXMax = xMax || xExtentMax;
  // we delay the tween if previousYMax < yMax
  let yMaxStore = writable(yExtentMax);
  let previousYMax = previousValueStore(yMaxStore);

  $: yMaxStore.set(yExtentMax);
  const timeRangeKey = writable(`${xMin}-${xMax}`);

  const previousTimeRangeKey = previousValueStore(timeRangeKey);

  // FIXME: move this function to utils.ts
  /** reset the keys to trigger animations on time range changes */
  let syncTimeRangeKey;
  $: {
    timeRangeKey.set(`${xMin}-${xMax}`);
    if ($previousTimeRangeKey !== $timeRangeKey) {
      if (syncTimeRangeKey) clearTimeout(syncTimeRangeKey);
      syncTimeRangeKey = setTimeout(() => {
        previousTimeRangeKey.set($timeRangeKey);
      }, 400);
    }
  }

  $: delay =
    $previousTimeRangeKey === $timeRangeKey && $previousYMax < yExtentMax
      ? 100
      : 0;

  function alwaysBetween(min, max, value) {
    // note: must work with dates
    if (value < min) return min;
    if (value > max) return max;
    return value;
  }

  function inBounds(min, max, value) {
    return value >= min && value <= max;
  }

  function resetScrub() {
    metricsExplorerStore.setSelectedScrubRange(metricViewName, undefined);
  }

  // $: if (scrubbing) {
  //   scrubEnd = alwaysBetween(internalXMin, internalXMax, mouseoverValue?.x);
  // }

  function startScrub(event) {
    console.log("start", event.detail);

    const scrubStartDate = getBisectedTimeFromCordinates(
      event.detail?.start?.x,
      $xScale,
      xAccessor,
      data,
      TIME_GRAIN[timeGrain].label
    );

    if (hasSubrangeSelected) {
      // check if we are scrubbing on the edges of scrub rect
      // if (
      //   scrubStartDate?.getTime() === scrubStart?.getTime() ||
      //   scrubStartDate?.getTime() === scrubEnd?.getTime()
      // ) {
      //   metricsExplorerStore.setSelectedScrubRange(metricViewName, {
      //     ...$dashboardStore?.selectedScrubRange,
      //     isScrubbing: true,
      //   });
      // }
    } else {
      metricsExplorerStore.setSelectedScrubRange(metricViewName, {
        name: TimeRangePreset.CUSTOM,
        start: scrubStartDate,
        end: undefined,
        isScrubbing: true,
      });
    }
  }

  function moveScrub(event) {
    const intermediateScrubVal = getBisectedTimeFromCordinates(
      event.detail?.stop?.x,
      $xScale,
      xAccessor,
      data,
      TIME_GRAIN[timeGrain].label
    );

    if (hasSubrangeSelected && !scrubbing) {
      if (intermediateScrubVal < scrubStart) {
        console.log("inside left");
        if (intermediateScrubVal?.getTime() !== scrubStart?.getTime()) {
          metricsExplorerStore.setSelectedScrubRange(metricViewName, {
            name: TimeRangePreset.CUSTOM,
            start: intermediateScrubVal,
            end: scrubEnd,
            isScrubbing: false,
          });
        }
      }
    } else {
      // Only make state changes when the bisected value changes
      if (intermediateScrubVal?.getTime() !== scrubEnd?.getTime()) {
        metricsExplorerStore.setSelectedScrubRange(metricViewName, {
          name: TimeRangePreset.CUSTOM,
          start: scrubStart,
          end: intermediateScrubVal,
          isScrubbing: true,
        });
      }
    }
  }

  function endScrub(event) {
    scrubEnd = getBisectedTimeFromCordinates(
      event.detail?.stop?.x,
      $xScale,
      xAccessor,
      data,
      TIME_GRAIN[timeGrain].label
    );

    justCreatedScrub = true;
    // reset justCreatedScrub after 100 milliseconds
    setTimeout(() => {
      justCreatedScrub = false;
    }, 100);

    metricsExplorerStore.setSelectedScrubRange(metricViewName, {
      ...$dashboardStore?.selectedScrubRange,
      end: scrubEnd,
      isScrubbing: false,
    });
  }

  async function onMouseClick() {
    // skip if still scrubbing
    if (justCreatedScrub || scrubbing) return;

    // skip if no scrub range selected
    if (!hasSubrangeSelected) return;

    const { start, end } = getOrderedStartEnd(scrubStart, scrubEnd);
    if (mouseoverValue?.x < start || mouseoverValue?.x > end) {
      // console.log(
      //   scrubbing,
      //   "mouseoverValue",
      //   mouseoverValue?.x,
      //   hasSubrangeSelected
      // );
      resetScrub();
    }
  }
</script>

<SimpleDataGraphic
  overflowHidden={false}
  yMin={internalYMin}
  yMax={internalYMax}
  shareYScale={false}
  yType="number"
  xType="date"
  {width}
  {height}
  top={4}
  left={0}
  right={50}
  bind:mouseoverValue
  bind:hovered
  let:config
  let:yScale
  customClass={scrubbing ? "cursor-ew-resize" : ""}
  yMinTweenProps={tweenProps}
  yMaxTweenProps={tweenProps}
  xMaxTweenProps={tweenProps}
  xMinTweenProps={tweenProps}
  on:click={() => onMouseClick()}
  on:scrub-start={(e) => startScrub(e)}
  on:scrub-move={(e) => moveScrub(e)}
  on:scrub-end={(e) => endScrub(e)}
>
  <Axis side="right" {numberKind} />
  <Grid />
  <Body>
    <!-- key on the time range itself to prevent weird tweening animations.
    We'll need to migrate this to a more robust solution once we've figured out
    the right way to "tile" together a time series with multiple pages of data.
    -->
    <!-- {#key $timeRangeKey} -->
    {#if showComparison}
      <g
        class="transition-opacity"
        class:opacity-80={mouseoverValue?.x}
        class:opacity-40={!mouseoverValue?.x}
      >
        <ChunkedLine
          area={false}
          lineColor={`hsl(217, 10%, 60%)`}
          delay={$timeRangeKey !== $previousTimeRangeKey ? 0 : delay}
          duration={$timeRangeKey !== $previousTimeRangeKey ? 0 : 200}
          {data}
          {xAccessor}
          yAccessor="comparison.{yAccessor}"
        />
      </g>
    {/if}
    <ChunkedLine
      lineColor={mainLineColor}
      {areaColor}
      delay={$timeRangeKey !== $previousTimeRangeKey ? 0 : delay}
      duration={$timeRangeKey !== $previousTimeRangeKey ? 0 : 200}
      {data}
      {xAccessor}
      {yAccessor}
    />
    <!-- {/key} -->
    <line
      x1={config.plotLeft}
      x2={config.plotLeft + config.plotRight}
      y1={yScale(0)}
      y2={yScale(0)}
      class="stroke-blue-200"
    />
  </Body>
  {#if !scrubbing && mouseoverValue?.x}
    <WithRoundToTimegrain
      strategy={TimeRoundingStrategy.PREVIOUS}
      value={mouseoverValue.x}
      {timeGrain}
      let:roundedValue
    >
      <WithBisector
        {data}
        callback={(d) => d[xAccessor]}
        value={roundedValue}
        let:point
      >
        {#if point && inBounds(internalXMin, internalXMax, point[xAccessor])}
          <g transition:fly|local={{ duration: 100, x: -4 }}>
            <text
              class="fill-gray-600"
              style:paint-order="stroke"
              stroke="white"
              stroke-width="3px"
              x={config.plotLeft + config.bodyBuffer + 6}
              y={config.plotTop + 10 + config.bodyBuffer}
            >
              {mouseoverTimeFormat(point[labelAccessor])}
            </text>
            {#if showComparison && point[`comparison.${labelAccessor}`]}
              <text
                style:paint-order="stroke"
                stroke="white"
                stroke-width="3px"
                class="fill-gray-400"
                x={config.plotLeft + config.bodyBuffer + 6}
                y={config.plotTop + 24 + config.bodyBuffer}
              >
                {mouseoverTimeFormat(point[`comparison.${labelAccessor}`])} prev.
              </text>
            {/if}
          </g>
          <g transition:fly|local={{ duration: 100, x: -4 }}>
            <MeasureValueMouseover
              {point}
              {xAccessor}
              {yAccessor}
              {showComparison}
              {mouseoverFormat}
              {numberKind}
            />
          </g>
        {/if}
      </WithBisector>
    </WithRoundToTimegrain>
  {/if}
  {#if hasSubrangeSelected}
    <MeasureScrub
      start={scrubStart}
      stop={scrubEnd}
      timeGrainLabel={TIME_GRAIN[timeGrain].label}
      {data}
      {xAccessor}
      {yAccessor}
      {mouseoverTimeFormat}
    />
  {/if}
</SimpleDataGraphic>
