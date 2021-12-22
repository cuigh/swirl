import { nextTick } from 'vue'
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { LoadingBarApi } from 'naive-ui'
import ForbiddenPage from '../pages/403.vue'
import NotFoundPage from '../pages/404.vue'
import LoginPage from '../pages/Login.vue'
import InitPage from '../pages/Init.vue'
import { store } from "../store";
import { t } from "@/locales";

var loadingBar: LoadingBarApi;

export function initLoadingBar(bar: LoadingBarApi) {
  loadingBar = bar
}

export function go(name: string, params: any) {
  router.push({ name: name, params: params })
}

const routes: RouteRecordRaw[] = [
  {
    name: 'home',
    path: "/",
    component: () => import('../pages/Home.vue'),
    meta: {
      auth: '?',
    }
  },
  {
    name: 'login',
    path: '/login',
    component: LoginPage,
    meta: {
      layout: "empty",
      auth: '*',
    }
  },
  {
    name: 'init',
    path: '/init',
    component: InitPage,
    meta: {
      layout: "empty",
      auth: '*',
    }
  },
  {
    name: 'profile',
    path: "/profile",
    component: () => import('../pages/Profile.vue'),
    meta: {
      auth: '?',
    }
  },
  {
    name: 'node_list',
    path: "/swarm/nodes",
    component: () => import('../pages/node/List.vue'),
    meta: {
      auth: 'node.view',
    }
  },
  {
    name: 'node_detail',
    path: "/swarm/nodes/:id",
    component: () => import('../pages/node/View.vue'),
    meta: {
      auth: 'node.view',
    }
  },
  {
    name: 'node_edit',
    path: "/swarm/nodes/:id/edit",
    component: () => import('../pages/node/Edit.vue'),
    meta: {
      auth: 'node.edit',
    }
  },
  {
    name: 'registry_list',
    path: "/swarm/registries",
    component: () => import('../pages/registry/List.vue'),
    meta: {
      auth: 'registry.view',
    }
  },
  {
    name: 'registry_detail',
    path: "/swarm/registries/:id",
    component: () => import('../pages/registry/View.vue'),
    meta: {
      auth: 'registry.view',
    }
  },
  {
    name: 'registry_new',
    path: "/swarm/registries/new",
    component: () => import('../pages/registry/Edit.vue'),
    meta: {
      auth: 'registry.edit',
    }
  },
  {
    name: 'registry_edit',
    path: "/swarm/registries/:id/edit",
    component: () => import('../pages/registry/Edit.vue'),
    meta: {
      auth: 'registry.edit',
    }
  },
  {
    name: 'network_list',
    path: "/swarm/networks",
    component: () => import('../pages/network/List.vue'),
    meta: {
      auth: 'network.view',
    }
  },
  {
    name: 'network_new',
    path: "/swarm/networks/new",
    component: () => import('../pages/network/New.vue'),
    meta: {
      auth: 'network.edit',
    }
  },
  {
    name: 'network_detail',
    path: "/swarm/networks/:name",
    component: () => import('../pages/network/View.vue'),
    meta: {
      auth: 'network.view',
    }
  },
  {
    name: "service_list",
    path: "/swarm/services",
    component: () => import('../pages/service/List.vue'),
    meta: {
      auth: 'service.view',
    }
  },
  {
    name: "service_detail",
    path: "/swarm/services/:name",
    component: () => import('../pages/service/View.vue'),
    meta: {
      auth: 'service.view',
    }
  },
  {
    name: "service_new",
    path: "/swarm/services/new",
    component: () => import('../pages/service/Edit.vue'),
    meta: {
      auth: 'service.edit',
    }
  },
  {
    name: "service_edit",
    path: "/swarm/services/:name/edit",
    component: () => import('../pages/service/Edit.vue'),
    meta: {
      auth: 'service.edit',
    }
  },
  {
    name: "task_list",
    path: "/swarm/tasks",
    component: () => import('../pages/task/List.vue'),
    meta: {
      auth: 'task.view',
    }
  },
  {
    name: "task_detail",
    path: "/swarm/tasks/:id",
    component: () => import('../pages/task/View.vue'),
    meta: {
      auth: 'task.view',
    }
  },
  {
    name: "config_list",
    path: "/swarm/configs",
    component: () => import('../pages/config/List.vue'),
    meta: {
      auth: 'config.view',
    }
  },
  {
    name: "config_detail",
    path: "/swarm/configs/:id",
    component: () => import('../pages/config/View.vue'),
    meta: {
      auth: 'config.view',
    }
  },
  {
    name: "config_new",
    path: "/swarm/configs/new",
    component: () => import('../pages/config/Edit.vue'),
    meta: {
      auth: 'config.edit',
    }
  },
  {
    name: "config_edit",
    path: "/swarm/configs/:id/edit",
    component: () => import('../pages/config/Edit.vue'),
    meta: {
      auth: 'config.edit',
    }
  },
  {
    name: "secret_list",
    path: "/swarm/secrets",
    component: () => import('../pages/secret/List.vue'),
    meta: {
      auth: 'secret.view',
    }
  },
  {
    name: "secret_detail",
    path: "/swarm/secrets/:id",
    component: () => import('../pages/secret/View.vue'),
    meta: {
      auth: 'secret.view',
    }
  },
  {
    name: "secret_new",
    path: "/swarm/secrets/new",
    component: () => import('../pages/secret/Edit.vue'),
    meta: {
      auth: 'secret.edit',
    }
  },
  {
    name: "secret_edit",
    path: "/swarm/secrets/:id/edit",
    component: () => import('../pages/secret/Edit.vue'),
    meta: {
      auth: 'secret.edit',
    }
  },
  {
    name: "stack_list",
    path: "/swarm/stacks",
    component: () => import('../pages/stack/List.vue'),
    meta: {
      auth: 'stack.view',
    }
  },
  {
    name: "stack_detail",
    path: "/swarm/stacks/:name",
    component: () => import('../pages/stack/View.vue'),
    meta: {
      auth: 'stack.view',
    }
  },
  {
    name: "stack_new",
    path: "/swarm/stacks/new",
    component: () => import('../pages/stack/Edit.vue'),
    meta: {
      auth: 'stack.edit',
    }
  },
  {
    name: "stack_edit",
    path: "/swarm/stacks/:name/edit",
    component: () => import('../pages/stack/Edit.vue'),
    meta: {
      auth: 'stack.edit',
    }
  },
  {
    name: "image_list",
    path: "/local/images",
    component: () => import('../pages/image/List.vue'),
    meta: {
      auth: 'image.view',
    }
  },
  {
    name: "image_detail",
    path: "/local/images/:node/:id",
    component: () => import('../pages/image/View.vue'),
    meta: {
      auth: 'image.view',
    }
  },
  {
    name: "container_list",
    path: "/local/containers",
    component: () => import('../pages/container/List.vue'),
    meta: {
      auth: 'container.view',
    }
  },
  {
    name: "container_detail",
    path: "/local/containers/:node/:id",
    component: () => import('../pages/container/View.vue'),
    meta: {
      auth: 'container.view',
    }
  },
  {
    name: "volume_list",
    path: "/local/volumes",
    component: () => import('../pages/volume/List.vue'),
    meta: {
      auth: 'volume.view',
    }
  },
  {
    name: "volume_detail",
    path: "/local/volumes/:node/:name",
    component: () => import('../pages/volume/View.vue'),
    meta: {
      auth: 'volume.view',
    }
  },
  {
    name: "volume_new",
    path: "/local/volumes/:node/new",
    component: () => import('../pages/volume/New.vue'),
    meta: {
      auth: 'volume.edit',
    }
  },
  {
    name: "user_list",
    path: "/system/users",
    component: () => import('../pages/user/List.vue'),
    meta: {
      auth: 'user.view',
    }
  },
  {
    name: "user_new",
    path: "/system/users/new",
    component: () => import('../pages/user/Edit.vue'),
    meta: {
      auth: 'user.edit',
    }
  },
  {
    name: "user_detail",
    path: "/system/users/:id",
    component: () => import('../pages/user/View.vue'),
    meta: {
      auth: 'user.view',
    }
  },
  {
    name: "user_edit",
    path: "/system/users/:id/edit",
    component: () => import('../pages/user/Edit.vue'),
    meta: {
      auth: 'user.edit',
    }
  },
  {
    name: "role_list",
    path: "/system/roles",
    component: () => import('../pages/role/List.vue'),
    meta: {
      auth: 'role.view',
    }
  },
  {
    name: "role_new",
    path: "/system/roles/new",
    component: () => import('../pages/role/Edit.vue'),
    meta: {
      auth: 'role.edit',
    }
  },
  {
    name: "role_detail",
    path: "/system/roles/:id",
    component: () => import('../pages/role/View.vue'),
    meta: {
      auth: 'role.view',
    }
  },
  {
    name: "role_edit",
    path: "/system/roles/:id/edit",
    component: () => import('../pages/role/Edit.vue'),
    meta: {
      auth: 'role.edit',
    }
  },
  {
    name: "event_list",
    path: "/system/events",
    component: () => import('../pages/event/List.vue'),
    meta: {
      auth: 'event.view',
    }
  },
  {
    name: "chart_list",
    path: "/system/charts",
    component: () => import('../pages/chart/List.vue'),
    meta: {
      auth: 'chart.view',
    }
  },
  {
    name: "chart_detail",
    path: "/system/charts/:id",
    component: () => import('../pages/chart/View.vue'),
    meta: {
      auth: 'chart.view',
    }
  },
  {
    name: "chart_new",
    path: "/system/charts/new",
    component: () => import('../pages/chart/Edit.vue'),
    meta: {
      auth: 'chart.edit',
    }
  },
  {
    name: "chart_edit",
    path: "/system/charts/:id/edit",
    component: () => import('../pages/chart/Edit.vue'),
    meta: {
      auth: 'chart.edit',
    }
  },
  {
    name: "setting",
    path: "/system/settings",
    component: () => import('../pages/setting/Setting.vue'),
    meta: {
      auth: 'setting.view',
    }
  },
  {
    name: '403',
    path: '/403',
    component: ForbiddenPage,
    meta: {
      layout: "simple",
      auth: '*',
    }
  },
  {
    name: '404',
    path: '/404',
    component: NotFoundPage,
    meta: {
      layout: "simple",
      auth: '*',
    }
  },
  {
    name: 'not-found',
    path: '/:pathMatch(.*)*',
    redirect: { name: '404' }
  },
]

function createSiteRouter() {
  const router = createRouter({
    history: createWebHistory(),
    routes,
  })

  router.beforeEach(function (to, from, next) {
    if (!from || to.path !== from.path) {
      loadingBar?.start()
      window.document.title = t(`titles.${to.name as string}`) + ' - Swirl'
    }

    const auth = to.meta.auth || '*'
    if (auth !== '*') {
      if (store.getters.anonymous) {
        next({ name: 'login', query: { redirect: to.fullPath } })
        return
      }

      if (auth !== '?' && !store.getters.allow(auth)) {
        next({ name: '403' })
        return
      }
    }

    next()
  })

  router.afterEach(function (to, from) {
    if (!from || to.path !== from.path) {
      loadingBar?.finish()
      if (to.hash && to.hash !== from.hash) {
        nextTick(() => {
          const el = document.querySelector(to.hash)
          if (el) el.scrollIntoView()
        })
      }
    }
  })

  return router
}

export const router = createSiteRouter()