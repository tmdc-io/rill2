export function fetchTokenFromLocalStorage() {
    const storeValue = localStorage.getItem("userInfo") ? JSON.parse(localStorage.getItem("userInfo")) : {};
    const token = storeValue?.accessToken;
    return token;
  }