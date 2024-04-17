<script lang="ts">
    const browserURL = new URL(window.location.href);
    const url = browserURL.origin+"/home";
</script>
<style>
    .default-text {

        padding:4px; margin: 4px; font-size: 16px !important; font-family: times; color: #000000; font-variation-settings: normal !important; 
        -webkit-font-smoothing: auto !important; 
    }
</style>
<div class="default-text">
    Hello stranger! We could not find your JWT Token (<a href="https://en.wikipedia.org/wiki/JSON_Web_Token" target="_blank">Learn more</a>) Please click here to <a href={url} target="_blank">Go Home</a>. It will automatically generate your token.
</div>