<script lang="ts">
  import { UserManager, UserManagerSettings, User } from 'oidc-client-ts';

// Define OIDC Configuration
 const oidcConfig: UserManagerSettings = {
    authority: process.env.VITE_OIDC_AUTHORITY,
    client_id: process.env.VITE_OIDC_CLIENT_ID,
    redirect_uri: process.env.VITE_OIDC_REDIRECT_URI,
    post_logout_redirect_uri: process.env.VITE_OIDC_POST_LOGOUT_REDIRECT_URI,
    response_type: process.env.VITE_OIDC_RESPONSE_TYPE || 'code',
    scope: process.env.VITE_OIDC_SCOPE || 'openid profile email',
    client_secret: process.env.VITE_OIDC_CLIENT_SECRET,
    loadUserInfo: process.env.VITE_OIDC_LOADUSERINFO,
    automaticSilentRenew: process.env.VITE_OIDC_AUTOMATIC_SILENT_RENEW,
    checkSessionIntervalInSeconds: process.env.VITE_OIDC_CHECK_SESSION_INTERVAL,
    revokeTokensOnSignout: process.env.VITE_OIDC_ACCESSTOKENEXPIRINGNOTIFICATIONTIME
  };

  const userManager = new UserManager(oidcConfig);

  userManager.signinSilentCallback().then(() => {
    console.log('Silent renew successful');
  }).catch((error) => {
    console.error('Silent renew error:', error);
  });
</script>

<p>Silent renew processing...</p>
