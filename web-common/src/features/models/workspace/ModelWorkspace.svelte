<script lang="ts">
  import { EntityType } from "@rilldata/web-common/features/entity-management/types";
  import { appStore } from "@rilldata/web-local/lib/application-state-stores/app-store";
  import { runtimeStore } from "@rilldata/web-local/lib/application-state-stores/application-store";
  import Tab from "../../../components/tab/Tab.svelte";
  import TabGroup from "../../../components/tab/TabGroup.svelte";
  import { WorkspaceContainer } from "../../../layout/workspace";
  import SQLAssistant from "../open-ai/SQLAssistant.svelte";
  import ModelInspector from "./inspector/ModelInspector.svelte";
  import ModelBody from "./ModelBody.svelte";
  import ModelWorkspaceHeader from "./ModelWorkspaceHeader.svelte";

  export let modelName: string;
  export let focusEditorOnMount = false;

  const switchToModel = async (modelName: string) => {
    if (!modelName) return;

    appStore.setActiveEntity(modelName, EntityType.Model);
  };

  $: switchToModel(modelName);

  $: isGpt3Enabled = $runtimeStore.openAIAPIKey !== "";
  let selectedInspectorTab: "profile" | "gpt" = "profile";
</script>

{#key modelName}
  <WorkspaceContainer assetID={modelName}>
    <div slot="header">
      <ModelWorkspaceHeader {modelName} />
    </div>
    <div slot="body">
      <ModelBody {modelName} {focusEditorOnMount} />
    </div>
    <div slot="inspector">
      {#if isGpt3Enabled}
        <div class="mx-2">
          <TabGroup
            variant="simple"
            on:select={(event) => {
              selectedInspectorTab = event.detail;
            }}
          >
            <Tab
              compact
              selected={selectedInspectorTab === "profile"}
              value={"profile"}>Profiler</Tab
            >
            <Tab compact selected={selectedInspectorTab === "gpt"} value={"gpt"}
              >SQL Assistant</Tab
            >
          </TabGroup>
        </div>
        <hr />
        {#if selectedInspectorTab === "profile"}
          <ModelInspector {modelName} />
        {:else if selectedInspectorTab === "gpt"}
          <SQLAssistant {modelName} />
        {/if}
      {:else}
        <ModelInspector {modelName} />
      {/if}
    </div>
  </WorkspaceContainer>
{/key}
