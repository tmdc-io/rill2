<script lang="ts">
    import CtaButton from "@rilldata/web-common/components/calls-to-action/CTAButton.svelte";
    import CtaContentContainer from "@rilldata/web-common/components/calls-to-action/CTAContentContainer.svelte";
    import CtaLayoutContainer from "@rilldata/web-common/components/calls-to-action/CTALayoutContainer.svelte";
    import CtaMessage from "@rilldata/web-common/components/calls-to-action/CTAMessage.svelte";

    export let statusCode = 401;
    export let header = "Unauthorized";
    export let body = "Oops! It seems you're logged out";
  
    function handleRedirection() {
        const browserURL = new URL(window.location.href);
        const origin = browserURL.origin;
        window.open(origin+"/home"); 
        body = "Refresh screen!";
    }
</script>

<CtaLayoutContainer>
    <CtaContentContainer>
      {#if statusCode}
        <h1
          class="text-8xl font-extrabold bg-gradient-to-b from-[#CBD5E1] to-[#E2E8F0] text-transparent bg-clip-text"
        >
          {statusCode}
        </h1>
      {/if}
      <h2 class="text-lg font-semibold">{header}</h2>
      <CtaMessage>
        {body}
      </CtaMessage>
      <CtaButton variant="primary-outline" on:click={() => handleRedirection()}>
        Go Home!
      </CtaButton>
    </CtaContentContainer>
</CtaLayoutContainer>
