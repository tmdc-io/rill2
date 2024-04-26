<script lang="ts">
  import { fly } from "svelte/transition";
  import BackgroundImage from "./BackgroundImage.svelte";
  import ProjectCards from "./ProjectCards.svelte";
  import TitleContent from "./TitleContent.svelte";
  import UserTestCta from "./UserTestCTA.svelte";
  import RadixH1 from "@rilldata/web-common/components/typography/RadixH1.svelte";
</script>

<!--<BackgroundImage>-->
  <div in:fly|global={{ duration: 1600, delay: 400, y: 8 }}>
    <div class="flex flex-col items-center gap-y-6 py-10 px-10 2xl:py-16">
      <RadixH1>
        <span class="bg-gradient-to-r from-primary-900 to-primary-800 text-transparent bg-clip-text opacity-75">
          Opps! You appear to have wandered off, Alice!
        </span>
      </RadixH1>

<!--      <TitleContent />-->
<!--      <div-->
<!--        class="flex flex-col items-center gap-y-6"-->
<!--        in:fly|global={{ duration: 1600, delay: 1200, y: 4 }}-->
<!--      >-->
<!--        <ProjectCards />-->
<!--        <UserTestCta />-->
<!--      </div>-->
    </div>
  </div>
<!--</BackgroundImage>-->
