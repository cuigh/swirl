import { h } from 'vue'
import { NIcon, MenuOption } from 'naive-ui'
import { RouteLocationNormalizedLoaded, RouterLink } from 'vue-router'
import {
  HomeOutline,
  PersonOutline,
  PeopleOutline,
  SettingsOutline,
  ConstructOutline,
  GridOutline,
  GlobeOutline,
  CubeOutline,
  BarChartOutline,
  LayersOutline,
  DocumentTextOutline,
  DocumentOutline,
  DocumentLockOutline,
  FileTrayFullOutline,
  BusinessOutline,
  ServerOutline,
  AlbumsOutline,
  ImageOutline,
  ImagesOutline,
} from "@vicons/ionicons5";
import XIcon from "@/components/Icon.vue";
import { t } from "@/locales";

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

export const renderMenuLabel = (option: any) => {
  if (!('path' in option)) {
    return option.label
  }
  return h(
    RouterLink,
    {
      to: option.path
    },
    {
      default: () => option.label
    }
  )
}

export function findMenuValue(route: RouteLocationNormalizedLoaded): string {
  var path = route.path;
  do {
    const option = findOption(menuOptions, path)
    if (option) {
      return option.key
    } else {
      const index = path.lastIndexOf("/")
      if (index <= 0) {
        return ""
      }
      path = path.substring(0, index)
    }
  } while (true)
}

function findOption(options: MenuOption[], path: string): any {
  for (const option of options) {
    if (option.path === path) {
      return option
    } else if (option.children) {
      const opt = findOption(option.children, path)
      if (opt) return opt
    }
  }
  return null
}

export function findActiveOptions(route: RouteLocationNormalizedLoaded): MenuOption[] {
  const result: MenuOption[] = []
  findOptions(result, menuOptions, route.path)
  return result
}

function findOptions(result: MenuOption[], options: MenuOption[], path: string): boolean {
  for (const option of options) {
    if (option.path) {
      if (option.path != "/" && path.startsWith(<string>option.path)) {
        result.push(option)
        return true
      }
    } else if (option.children) {
      result.push(option)
      if (findOptions(result, option.children, path)) {
        return true
      } else {
        result.pop()
      }
    }
  }
  return false
}

export const menuOptions: MenuOption[] = [
  {
    label: t('fields.home'),
    key: "home",
    path: "/",
    icon: renderIcon(HomeOutline),
  },
  {
    label: t('fields.swarm'),
    key: "swarm",
    icon: renderIcon(GridOutline),
    children: [
      {
        label: t('objects.registry'),
        key: "registries",
        path: "/swarm/registries",
        icon: renderIcon(BusinessOutline),
      },
      {
        label: t('objects.node'),
        key: "nodes",
        path: "/swarm/nodes",
        icon: renderIcon(ServerOutline),
      },
      {
        label: t('objects.network'),
        key: "networks",
        path: "/swarm/networks",
        icon: renderIcon(GlobeOutline),
      },
      {
        label: t('objects.service'),
        key: "services",
        path: "/swarm/services",
        icon: renderIcon(ImageOutline),
      },
      {
        label: t('objects.task'),
        key: "tasks",
        path: "/swarm/tasks",
        icon: renderIcon(ImagesOutline),
      },
      {
        label: t('objects.stack'),
        key: "stacks",
        path: "/swarm/stacks",
        icon: renderIcon(AlbumsOutline),
      },
      {
        label: t('objects.config'),
        key: "configs",
        path: "/swarm/configs",
        icon: renderIcon(DocumentOutline),
      },
      {
        label: t('objects.secret'),
        key: "secrets",
        path: "/swarm/secrets",
        icon: renderIcon(DocumentLockOutline),
      },
    ],
  },
  {
    label: t('fields.local'),
    key: "local",
    icon: renderIcon(CubeOutline),
    children: [
      {
        label: t('objects.image'),
        key: "images",
        path: "/local/images",
        icon: renderIcon(LayersOutline),
      },
      {
        label: t('objects.container'),
        key: "containers",
        path: "/local/containers",
        icon: () => h(XIcon, {
          path: [
            'M28 12h-8V4h8zm-6-2h4V6h-4z',
            'M17 15V9H9v14h14v-8zm-6-4h4v4h-4zm4 10h-4v-4h4zm6 0h-4v-4h4z',
            'M26 28H6a2.002 2.002 0 0 1-2-2V6a2.002 2.002 0 0 1 2-2h10v2H6v20h20V16h2v10a2.002 2.002 0 0 1-2 2z',
          ],
          viewBox: '0 0 32 32'
        }),
      },
      {
        label: t('objects.volume'),
        key: "volumes",
        path: "/local/volumes",
        icon: renderIcon(FileTrayFullOutline),
      },
    ],
  },
  {
    label: t('fields.system'),
    key: "system",
    icon: renderIcon(SettingsOutline),
    children: [
      {
        label: t('objects.user'),
        key: "users",
        path: "/system/users",
        icon: renderIcon(PersonOutline),
      },
      {
        label: t('objects.role'),
        key: "roles",
        path: "/system/roles",
        icon: renderIcon(PeopleOutline),
      },
      {
        label: t('objects.chart'),
        key: "charts",
        path: "/system/charts",
        icon: renderIcon(BarChartOutline),
      },
      {
        label: t('objects.event'),
        key: "events",
        path: "/system/events",
        icon: renderIcon(DocumentTextOutline),
      },
      {
        label: t('objects.setting'),
        key: "config",
        path: "/system/settings",
        icon: renderIcon(ConstructOutline),
      },
    ],
  },
]
