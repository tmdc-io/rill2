<!-- @component 
A container GraphicContext for the time series in a metrics dashboard.
-->
<script lang="ts">
  import { GraphicContext } from "@rilldata/web-common/components/data-graphic/elements";
  // import { MEASURE_CONFIG } from "../config";
  // export let start: Date;
  // export let end: Date;
  // export let workspaceWidth: number;
</script>

<div class="grid w-max">
  <GraphicContext>
    <slot />
  </GraphicContext>
</div>
