export async function register() {
  if (typeof globalThis.localStorage !== "undefined" && typeof window === "undefined") {
    // Node.js 22+ exposes a broken localStorage global that crashes SSR
    // @ts-expect-error intentionally removing the broken global
    delete globalThis.localStorage;
  }
}
