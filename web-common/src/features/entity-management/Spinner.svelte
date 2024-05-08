<script lang="ts">
    import {EntityStatus} from "@rilldata/web-common/features/entity-management/types";

    export let size = "1em";
    export let status: EntityStatus = EntityStatus.Idle;

    export let duration = 500;
</script>

<div
  class="status"
  class:running={status === EntityStatus.Running}
  class:idle={status === EntityStatus.Idle}
  style="
    --status-transition: {duration}ms;
    --size: {size};
    width: {size};
    height: {size};"
/>

<style>
    .running {
        border: 2px solid #FFF;
        border-bottom-color: #d2d2d2;
        border-radius: 50%;
        display: inline-block;
        box-sizing: border-box;
        animation: rotation 1s linear infinite;
    }

    .idle {
        border: 2px solid #FFF;
        border-color: #d2d2d2;
        border-radius: 50%;
        display: inline-block;
        box-sizing: border-box;
        animation: rotation 1s linear infinite;
    }

    @keyframes rotation {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
</style>