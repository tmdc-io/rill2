<script lang="ts">
  import { notifications } from "@rilldata/web-common/components/notifications";
  import Tooltip from "@rilldata/web-common/components/tooltip/Tooltip.svelte";
  import TooltipContent from "@rilldata/web-common/components/tooltip/TooltipContent.svelte";
  import { forecastStore } from "@rilldata/web-local/lib/application-state-stores/app-store";
  import { createForm } from "svelte-forms-lib";
  import type { Readable } from "svelte/store";
  import type { MetricsInternalRepresentation } from "../../metrics-internal-store";
  import {
    CONFIG_SELECTOR,
    CONFIG_TOP_LEVEL_LABEL_CLASSES,
    INPUT_ELEMENT_CONTAINER,
    SELECTOR_BUTTON_TEXT_CLASSES,
    SELECTOR_CONTAINER,
  } from "../styles";

  $: period = $forecastStore;

  $: currentPeriod = period;

  const { form, handleSubmit } = createForm({
    initialValues: {
      period: period,
    },
    onSubmit: async (value) => {
      try {
        forecastStore.set(value.period);
      } catch (err) {
        console.error(err);
        notifications.send({ message: err.response.data.message });
      }
    },
  });

  function handleKeydown(event: KeyboardEvent) {
    if (event.code == "Enter") {
      event.preventDefault();
      handleSubmit(event);
      (event.target as HTMLInputElement).blur();
    }
  }

  function updateFormWithNewDisplayName(period: string) {
    $form.period = period;
  }

  // This kicks in when the user changes the display name via code artifact
  $: updateFormWithNewDisplayName(currentPeriod);
</script>

<div
  class={INPUT_ELEMENT_CONTAINER.classes}
  style={INPUT_ELEMENT_CONTAINER.style}
>
  <Tooltip alignment="middle" distance={8} location="bottom">
    <div class={CONFIG_TOP_LEVEL_LABEL_CLASSES}>Forecast Period</div>

    <TooltipContent slot="tooltip-content">
      Mention the number of periods you want to support
    </TooltipContent>
  </Tooltip>
  <div class={SELECTOR_CONTAINER.classes} style={SELECTOR_CONTAINER.style}>
    <form id="display-name-form" autocomplete="off">
      <input
        type="text"
        bind:value={$form["period"]}
        on:keydown={handleKeydown}
        on:blur={handleSubmit}
        placeholder={"3"}
        class="{SELECTOR_BUTTON_TEXT_CLASSES} placeholder:font-normal placeholder:text-gray-600 font-semibold bg-white w-full hover:bg-gray-200 rounded border border-6 border-gray-200 hover:border-gray-300  hover:text-gray-900 px-2 py-1 h-[34px] {CONFIG_SELECTOR.focus}"
      />
    </form>
  </div>
</div>
