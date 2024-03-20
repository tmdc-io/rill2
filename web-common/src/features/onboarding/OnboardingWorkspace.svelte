<script lang="ts">
  import { IconSpaceFixer } from "../../components/button";
  import Button from "../../components/button/Button.svelte";
  import Add from "../../components/icons/Add.svelte";
  import { WorkspaceContainer } from "../../layout/workspace";
  import { createRuntimeServiceGetInstance } from "../../runtime-client";
  import { runtime } from "../../runtime-client/runtime-store";
  import { addSourceModal } from "../sources/modal/add-source-visibility";

  let steps: OnboardingStep[];
  $: instance = createRuntimeServiceGetInstance($runtime.instanceId);
  $: olapConnector = $instance.data?.instance?.olapConnector;
  $: if (olapConnector) {
    steps = olapConnector === "duckdb" ? duckDbSteps : nonDuckDbSteps;
  }

  interface OnboardingStep {
    id: string;
    heading: string;
    description: string;
  }

  // Onboarding steps for DuckDB OLAP driver
  const duckDbSteps: OnboardingStep[] = [
    {
      id: "source",
      heading: "Import your data source",
      description:
        "Click 'Add data' or drag a file (Parquet, NDJSON, or CSV) into this window.",
    },
    {
      id: "model",
      heading: "Model your sources into one big table",
      description:
        "Build intuition about your sources and use SQL to model them into an analytics-ready resource.",
    },
    {
      id: "metrics",
      heading: "Define your metrics and dimensions",
      description:
        "Define aggregate metrics and break out dimensions for your modeled data.",
    },
    {
      id: "dashboard",
      heading: "Explore your metrics dashboard",
      description:
        "Interactively explore line charts and leaderboards to uncover insights.",
    },
  ];

  // Onboarding steps for non-DuckDB OLAP drivers (ClickHouse, Druid)
  const nonDuckDbSteps: OnboardingStep[] = [
    {
      id: "table",
      heading: "Explore your tables",
      description:
        "Find your database tables in the left-hand-side navigational sidebar.",
    },
    {
      id: "metrics",
      heading: "Define your metrics and dimensions",
      description:
        "Define aggregate metrics and break out dimensions for your tables.",
    },
    {
      id: "dashboard",
      heading: "Explore your metrics dashboard",
      description:
        "Interactively explore line charts and leaderboards to uncover insights.",
    },
  ];
</script>

<WorkspaceContainer inspector={false}>
  <div class="pt-20 px-8 flex flex-col gap-y-6 items-center" slot="body">
    <div class="text-center">
      <div class="font-bold">Lens2 Infoboard</div>
      <p>Building data intuition at every step of analysis</p>
    </div>
  </div>
</WorkspaceContainer>
