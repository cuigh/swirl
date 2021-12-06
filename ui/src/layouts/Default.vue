<template>
  <n-layout :position="isMobile ? 'static' : 'absolute'">
    <n-layout-header bordered>
      <div class="header-left" align="center">
        <n-popover
          v-if="isMobile || isTablet"
          style="padding: 0; width: 200px"
          placement="bottom-end"
          display-directive="show"
          trigger="click"
          ref="menuPopover"
        >
          <template #trigger>
            <n-button size="small" style="margin-right: 8px">
              <template #icon>
                <n-icon>
                  <menu-outline />
                </n-icon>
              </template>
            </n-button>
          </template>
          <div style="overflow: auto; max-height: 79vh">
            <n-menu
              :value="menuValue"
              :options="menuOptions"
              :indent="18"
              @update:value="menuPopover.setShow(false)"
              :render-label="renderMenuLabel"
            />
          </div>
        </n-popover>
        <n-text tag="div" class="logo" :depth="1" @click="$router.push('/')">
          <img src="/favicon.ico" v-if="!isMobile" />
          Swirl
        </n-text>
      </div>
      <n-space justify="end" align="center" class="header-right" :size="0">
        <div style="margin-right: 10px; line-height: 56px">
          <n-text depth="3">v{{ version.version }}</n-text>
        </div>
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button
              type="default"
              size="small"
              :bordered="false"
              tag="a"
              href="https://github.com/cuigh/swirl"
              target="_blank"
            >
              <template #icon>
                <n-icon>
                  <LogoGithub />
                </n-icon>
              </template>
            </n-button>
          </template>
          GitHub
        </n-tooltip>
        <n-dropdown @select="selectOption" trigger="hover" :options="dropdownOptions" show-arrow>
          <n-button quaternary size="small">
            <template #icon>
              <n-icon>
                <PersonOutline />
              </n-icon>
            </template>
            {{ store.state.name }}
          </n-button>
        </n-dropdown>
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button size="small" quaternary @click="logout">
              <template #icon>
                <n-icon>
                  <LogOutOutline />
                </n-icon>
              </template>
            </n-button>
          </template>
          {{ t('buttons.sign_out') }}
        </n-tooltip>
      </n-space>
    </n-layout-header>
    <n-layout
      has-sider
      :position="isMobile ? 'static' : 'absolute'"
      :style="isMobile ? '' : 'top: 56px; bottom: 64px'"
    >
      <n-layout-sider
        v-if="!isMobile && !isTablet"
        bordered
        width="200"
        :collapsed-width="64"
        :collapsed="collapsed"
        collapse-mode="width"
        show-trigger="bar"
        trigger-style="right: -25px"
        collapsed-trigger-style="right: -25px"
        @collapse="collapsed = true"
        @expand="collapsed = false"
      >
        <n-menu
          :value="menuValue"
          :options="menuOptions"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :root-indent="20"
          :indent="24"
          :render-label="renderMenuLabel"
          :expanded-keys="expandedKeys"
          @update:expanded-keys="updateExpandedKeys"
        />
      </n-layout-sider>
      <n-layout-content>
        <router-view></router-view>
        <n-back-top :right="16" :bottom="10" />
      </n-layout-content>
    </n-layout>
    <n-layout-footer bordered :position="isMobile ? 'static' : 'absolute'">
      <span>{{ t('copyright') }}</span>
    </n-layout-footer>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted } from "vue";
import {
  NButton,
  NIcon,
  NMenu,
  NText,
  NSpace,
  NLayout,
  NLayoutHeader,
  NLayoutSider,
  NLayoutContent,
  NLayoutFooter,
  NPopover,
  NTooltip,
  NDropdown,
  NSwitch,
  NBackTop,
} from "naive-ui";
import { MenuOutline, PersonOutline, LogOutOutline, LogoGithub } from "@vicons/ionicons5";
import { RouterView, useRouter, useRoute } from "vue-router";
import { useStore } from "vuex";
import { useIsMobile, useIsTablet } from "@/utils";
import { findMenuValue, renderMenuLabel, menuOptions, findActiveOptions } from "@/router/menu";
import systemApi from "@/api/system";
import type { Version } from "@/api/system";
import { Mutations } from "@/store/mutations";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const dropdownOptions = [
  {
    label: t('titles.profile'),
    key: "profile",
  },
];
const store = useStore();
const router = useRouter();
const route = useRoute();
const menuPopover = ref();
const collapsed = ref(false)
const expandedKeys = ref([] as string[]);
const isMobile = useIsMobile()
const isTablet = useIsTablet()
const darkTheme = computed(() => store.state.preference.theme === "dark")
const menuValue = computed(() => findMenuValue(route))
const version = ref({} as Version);

function updateExpandedKeys(data: any) {
  expandedKeys.value = data
}

function selectOption(key: any) {
  switch (key as string) {
    case "profile":
      router.push("/profile")
      return
    default:
      console.info(key)
  }
}

function logout() {
  store.commit(Mutations.Logout);
  router.push("/login");
}

watch(() => route.path, (path: string) => {
  let keys = findActiveOptions(route).map((opt: any) => opt.key) as string[]
  expandedKeys.value = keys;
})

onMounted(async () => {
  const r = await systemApi.version();
  version.value = r.data as Version;
})
</script>

<style scoped>
::v-deep(.header-right .n-button__content) {
  margin-top: 4px;
}
.header-left {
  flex-grow: 1;
  width: 180px;
  display: flex;
  align-items: center;
}
.header-right {
  width: 320px;
}
/* .n-layout-header {
  background-color: #363636;
} */
.n-layout-sider {
  box-shadow: 2px 0 4px -2px rgb(10 10 10 / 10%);
}
.n-layout-footer {
  box-shadow: 0px -2px 4px -2px rgb(10 10 10 / 10%);
  /* background-image: radial-gradient(circle at 1% 1%,#328bf2,#1644ad); */
}
/* .n-layout-header {
  background-image: linear-gradient(to right, rgb(91, 121, 162) 0%, rgb(46, 68, 105) 100%);
}
.logo {
  color: white;
}
.n-layout-header .n-icon {
  color: white;
}
::v-deep(.n-layout-header .n-button__content) {
  color: white;
}
::v-deep(.n-layout-header .n-button__content:hover) {
  color: green;
} */
</style>
