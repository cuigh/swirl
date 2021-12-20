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
  },
  {
    name: 'login',
    path: '/login',
    component: LoginPage,
    meta: {
      layout: "empty",
      anonymous: true,
    }
  },
  {
    name: 'init',
    path: '/init',
    component: InitPage,
    meta: {
      layout: "empty",
      anonymous: true,
    }
  },
  {
    name: 'profile',
    path: "/profile",
    component: () => import('../pages/Profile.vue'),
  },
  {
    name: 'node_list',
    path: "/swarm/nodes",
    component: () => import('../pages/node/List.vue'),
  },
  {
    name: 'node_detail',
    path: "/swarm/nodes/:id",
    component: () => import('../pages/node/View.vue'),
  },
  {
    name: 'node_edit',
    path: "/swarm/nodes/:id/edit",
    component: () => import('../pages/node/Edit.vue'),
  },
  {
    name: 'registry_list',
    path: "/swarm/registries",
    component: () => import('../pages/registry/List.vue'),
  },
  {
    name: 'registry_detail',
    path: "/swarm/registries/:id",
    component: () => import('../pages/registry/View.vue'),
  },
  {
    name: 'registry_new',
    path: "/swarm/registries/new",
    component: () => import('../pages/registry/Edit.vue'),
  },
  {
    name: 'registry_edit',
    path: "/swarm/registries/:id/edit",
    component: () => import('../pages/registry/Edit.vue'),
  },
  {
    name: 'network_list',
    path: "/swarm/networks",
    component: () => import('../pages/network/List.vue'),
  },
  {
    name: 'network_new',
    path: "/swarm/networks/new",
    component: () => import('../pages/network/New.vue'),
  },
  {
    name: 'network_detail',
    path: "/swarm/networks/:name",
    component: () => import('../pages/network/View.vue'),
  },
  {
    name: "service_list",
    path: "/swarm/services",
    component: () => import('../pages/service/List.vue'),
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
  },
  {
    name: "task_detail",
    path: "/swarm/tasks/:id",
    component: () => import('../pages/task/View.vue'),
  },
  {
    name: "config_list",
    path: "/swarm/configs",
    component: () => import('../pages/config/List.vue'),
  },
  {
    name: "config_detail",
    path: "/swarm/configs/:id",
    component: () => import('../pages/config/View.vue'),
  },
  {
    name: "config_new",
    path: "/swarm/configs/new",
    component: () => import('../pages/config/Edit.vue'),
  },
  {
    name: "config_edit",
    path: "/swarm/configs/:id/edit",
    component: () => import('../pages/config/Edit.vue'),
  },
  {
    name: "secret_list",
    path: "/swarm/secrets",
    component: () => import('../pages/secret/List.vue'),
  },
  {
    name: "secret_detail",
    path: "/swarm/secrets/:id",
    component: () => import('../pages/secret/View.vue'),
  },
  {
    name: "secret_new",
    path: "/swarm/secrets/new",
    component: () => import('../pages/secret/Edit.vue'),
  },
  {
    name: "secret_edit",
    path: "/swarm/secrets/:id/edit",
    component: () => import('../pages/secret/Edit.vue'),
  },
  {
    name: "stack_list",
    path: "/swarm/stacks",
    component: () => import('../pages/stack/List.vue'),
  },
  {
    name: "stack_detail",
    path: "/swarm/stacks/:name",
    component: () => import('../pages/stack/View.vue'),
  },
  {
    name: "stack_new",
    path: "/swarm/stacks/new",
    component: () => import('../pages/stack/Edit.vue'),
  },
  {
    name: "stack_edit",
    path: "/swarm/stacks/:name/edit",
    component: () => import('../pages/stack/Edit.vue'),
  },
  {
    name: "image_list",
    path: "/local/images",
    component: () => import('../pages/image/List.vue'),
  },
  {
    name: "image_detail",
    path: "/local/images/:node/:id",
    component: () => import('../pages/image/View.vue'),
  },
  {
    name: "container_list",
    path: "/local/containers",
    component: () => import('../pages/container/List.vue'),
  },
  {
    name: "container_detail",
    path: "/local/containers/:node/:id",
    component: () => import('../pages/container/View.vue'),
  },
  {
    name: "volume_list",
    path: "/local/volumes",
    component: () => import('../pages/volume/List.vue'),
  },
  {
    name: "volume_detail",
    path: "/local/volumes/:node/:name",
    component: () => import('../pages/volume/View.vue'),
  },
  {
    name: "volume_new",
    path: "/local/volumes/:node/new",
    component: () => import('../pages/volume/New.vue'),
  },
  {
    name: "user_list",
    path: "/system/users",
    component: () => import('../pages/user/List.vue'),
  },
  {
    name: "user_new",
    path: "/system/users/new",
    component: () => import('../pages/user/Edit.vue'),
  },
  {
    name: "user_detail",
    path: "/system/users/:id",
    component: () => import('../pages/user/View.vue'),
  },
  {
    name: "user_edit",
    path: "/system/users/:id/edit",
    component: () => import('../pages/user/Edit.vue'),
  },
  {
    name: "role_list",
    path: "/system/roles",
    component: () => import('../pages/role/List.vue'),
  },
  {
    name: "role_new",
    path: "/system/roles/new",
    component: () => import('../pages/role/Edit.vue'),
  },
  {
    name: "role_detail",
    path: "/system/roles/:id",
    component: () => import('../pages/role/View.vue'),
  },
  {
    name: "role_edit",
    path: "/system/roles/:id/edit",
    component: () => import('../pages/role/Edit.vue'),
  },
  {
    name: "event_list",
    path: "/system/events",
    component: () => import('../pages/event/List.vue'),
  },
  {
    name: "chart_list",
    path: "/system/charts",
    component: () => import('../pages/chart/List.vue'),
  },
  {
    name: "chart_detail",
    path: "/system/charts/:id",
    component: () => import('../pages/chart/View.vue'),
  },
  {
    name: "chart_new",
    path: "/system/charts/new",
    component: () => import('../pages/chart/Edit.vue'),
  },
  {
    name: "chart_edit",
    path: "/system/charts/:id/edit",
    component: () => import('../pages/chart/Edit.vue'),
  },
  {
    name: "setting",
    path: "/system/settings",
    component: () => import('../pages/setting/Setting.vue'),
  },
  {
    name: '403',
    path: '/403',
    component: ForbiddenPage,
    meta: {
      layout: "simple",
      anonymous: true,
    }
  },
  {
    name: '404',
    path: '/404',
    component: NotFoundPage,
    meta: {
      layout: "simple",
      anonymous: true,
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

    if (to.matched.some(record => !record.meta.anonymous)) {
      // this route requires auth, if not logged in, redirect to login page.      
      if (store.getters.anonymous) {
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
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