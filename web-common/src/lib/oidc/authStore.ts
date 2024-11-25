import { writable } from 'svelte/store';
import { UserManager, UserManagerSettings, User } from 'oidc-client-ts';

// Define OIDC Configuration
const oidcConfig: UserManagerSettings = {
  authority: 'https://unique-haddock.dataos.app/oidc', // Replace with your OIDC provider URL
  client_id: 'dataos_generic', // Replace with your client ID
  redirect_uri: 'http://localhost:3000/dev/auth/callback', // Replace with your redirect URI
  post_logout_redirect_uri: 'http://localhost:3000/', // Replace with your logout redirect URI
  response_type: 'code', // OIDC response type
  scope: 'openid profile groups email federated:id', // Define required scopes
  client_secret: "665CAF977DA1ED77456B7CD3F3FD4",
  loadUserInfo: true
};

// Initialize UserManager
const userManager = new UserManager(oidcConfig);

// Define a writable store for managing the user's authentication state
export const user = writable<User | null>(null);

// Login function to initiate authentication
export async function login(): Promise<void> {
    console.log("clcik on login")
  try {
    await userManager.signinRedirect();
  } catch (error) {
    console.error('Error during login:', error);
  }
}

// Function to handle the OIDC callback and set the user
export async function handleCallback(): Promise<void> {
  try {
    const userResult = await userManager.signinRedirectCallback();
    console.log("userResult======>", userResult)
    localStorage.setItem('modern-oidc.user:https://unique-haddock.dataos.app/oidc:dataos_generic', JSON.stringify(userResult));
    user.set(userResult);
  } catch (error) {
    console.error('Error during callback handling:', error);
  }
}

// Logout function to end the session
export async function logout(): Promise<void> {
  try {
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
