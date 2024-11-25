<script lang="ts">
  import { onMount } from 'svelte';
  import { handleCallback } from "@rilldata/web-common/lib/oidc/authStore";

  onMount(async () => {
    console.log('Callback route accessed'); // Log to ensure this file is loaded
    try {
      await handleCallback();
      console.log('Callback handled successfully');
      window.location.href = '/'; // Redirect after successful handling
    } catch (error) {
      console.error('Error during callback:', error); // Log errors if they occur
      //window.location.href = '/error'; // Redirect to an error page
    }
  });
</script>

<p>Processing login...</p>
