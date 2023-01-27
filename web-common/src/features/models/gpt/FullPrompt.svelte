<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Button } from "../../../components/button";
  import { Dialog } from "../../../components/modal";
  import DialogFooter from "../../../components/modal/dialog/DialogFooter.svelte";

  export let prompt: string;

  const dispatch = createEventDispatcher();
</script>

<Dialog on:cancel={() => dispatch("close")} size="md" compact>
  <svelte:fragment slot="title">Full prompt for GPT3</svelte:fragment>
  <div class="m-4 flex-grow">
    <div>
      The following prompt is used for Open AI's <code>/completions</code> API
      with the <em>code-davinci-002</em> model:
    </div>
    <br />
    <div class="whitespace-pre overflow-auto">
      {prompt}
    </div>
  </div>

  <div class="bg-gray-100 border-t border-gray-300">
    <DialogFooter>
      <div class="flex items-center space-x-2">
        <Button type="primary" on:click={() => dispatch("close")}>Close</Button>
      </div>
    </DialogFooter>
  </div>
</Dialog>
