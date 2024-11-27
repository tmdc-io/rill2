import { writable } from 'svelte/store';
import { UserManager, UserManagerSettings, User } from 'oidc-client-ts';

// Define OIDC Configuration
const oidcConfig: UserManagerSettings = {
  authority: import.meta.env.VITE_OIDC_AUTHORITY,
  client_id: import.meta.env.VITE_OIDC_CLIENT_ID,
  redirect_uri: import.meta.env.VITE_OIDC_REDIRECT_URI,
  post_logout_redirect_uri: import.meta.env.VITE_OIDC_POST_LOGOUT_REDIRECT_URI,
  response_type: import.meta.env.VITE_OIDC_RESPONSE_TYPE || 'code',
  scope: import.meta.env.VITE_OIDC_SCOPE || 'openid profile email',
  client_secret: import.meta.env.VITE_OIDC_CLIENT_SECRET,
  loadUserInfo: import.meta.env.VITE_OIDC_LOADUSERINFO === 'true',
  automaticSilentRenew: import.meta.env.VITE_OIDC_AUTOMATIC_SILENT_RENEW === 'true',
  silent_redirect_uri: import.meta.env.VITE_OIDC_SILENT_REDIRECT_URI
};

console.log('OIDC Config:', oidcConfig);

// Initialize UserManager
const userManager = new UserManager(oidcConfig);

// Define a writable store for managing the user's authentication state
export const user = writable<User | null>(null);

// **Silent Renew Handlers**
userManager.events.addAccessTokenExpired(async () => {
  console.warn('Access token expired. Attempting silent renew...');
  try {
    await userManager.signinSilent();
    const renewedUser = await userManager.getUser();
    user.set(renewedUser);
    console.log('Silent renew successful:', renewedUser);
  } catch (error) {
    console.error('Silent renew failed:', error);
  }
});

userManager.events.addSilentRenewError((error) => {
  console.error('Silent renew error:', error);
});

// Login function to initiate authentication
export async function login(): Promise<void> {
  try {
    console.log('Initiating login...');
    await userManager.signinRedirect();
  } catch (error) {
    console.error('Error during login:', error);
  }
}

// Function to handle the OIDC callback and set the user
export async function handleCallback(): Promise<void> {
  try {
    const userResult = await userManager.signinRedirectCallback();
    console.log('User authenticated:', userResult);
    localStorage.setItem(
      'modern-oidc.user:https://unique-haddock.dataos.app/oidc:dataos_generic',
      JSON.stringify(userResult)
    );
    user.set(userResult);
  } catch (error) {
    console.error('Error during callback handling:', error);
  }
}

// Logout function to end the session
export async function logout(): Promise<void> {
  try {
    console.log('Logging out...');
    await userManager.signoutRedirect();
    user.set(null);
  } catch (error) {
    console.error('Error during logout:', error);
  }
}

// Utility function to fetch authenticated user's token
export async function getAccessToken(): Promise<string | null> {
  try {
    const currentUser = await userManager.getUser();
    return currentUser?.access_token || null;
  } catch (error) {
    console.error('Error fetching access token:', error);
    return null;
  }
}

// **Initialize Authentication**
export async function initializeAuth(): Promise<void> {
  try {
    const currentUser = await userManager.getUser();
    if (currentUser && !currentUser.expired) {
      console.log('User already logged in:', currentUser);
      user.set(currentUser);
    } else {
      console.log('No valid user session found.');
      user.set(null);
    }
  } catch (error) {
    console.error('Error during authentication initialization:', error);
  }
}
