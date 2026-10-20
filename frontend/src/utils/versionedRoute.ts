import { defineComponent, h, defineAsyncComponent, type AsyncComponentLoader } from 'vue'
import { useUiVersionStore } from '@/stores/uiVersion.store'

/**
 * Creates a reactive versioned route component that automatically switches between
 * V1 and V2 views based on `uiVersionStore.isV2`.
 *
 * If V2 is active but not provided (or during gradual rollout), it safely falls back to V1.
 */
export function createVersionedRoute(
  v1Loader: AsyncComponentLoader,
  v2Loader?: AsyncComponentLoader
) {
  const V1 = defineAsyncComponent(v1Loader)
  const V2 = v2Loader ? defineAsyncComponent(v2Loader) : V1

  return defineComponent({
    name: 'VersionedRouteWrapper',
    setup() {
      const uiVersionStore = useUiVersionStore()
      return () => {
        const ComponentToRender = (uiVersionStore.isV2 && v2Loader) ? V2 : V1
        return h(ComponentToRender)
      }
    }
  })
}

