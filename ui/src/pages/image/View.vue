<template>
  <x-page-header :subtitle="image.id">
    <template #action>
      <n-button secondary size="small" @click="$router.push({ name: 'image_list' })">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
    </template>
  </x-page-header>
  <div class="page-body">
    <n-tabs type="line" style="margin-top: -12px">
      <n-tab-pane name="detail" :tab="t('fields.detail')">
        <n-space vertical :size="16">
          <x-description label-placement="left" label-align="right" :label-width="110">
            <x-description-item :label="t('fields.id')" :span="2">{{ image.id }}</x-description-item>
            <x-description-item
              :label="t('fields.tags')"
              :span="2"
              v-if="image.tags && image.tags.length"
            >
              <n-space :size="4">
                <n-tag round size="small" type="default" v-for="tag in image.tags">{{ tag }}</n-tag>
              </n-space>
            </x-description-item>
            <x-description-item :label="t('fields.created_at')" :span="2">{{ image.created }}</x-description-item>
            <x-description-item :label="t('fields.size')">{{ formatSize(image.size) }}</x-description-item>
            <x-description-item :label="t('fields.platform')">{{ image.os + "/" + image.arch }}</x-description-item>
            <x-description-item
              :label="t('fields.docker_version')"
              v-if="image.dockerVersion"
              :span="2"
            >{{ image.dockerVersion }}</x-description-item>
            <x-description-item
              :label="t('fields.graph_driver')"
              v-if="image.graphDriver?.name"
            >{{ image.graphDriver?.name }}</x-description-item>
            <x-description-item
              :label="t('fields.root_fs')"
              v-if="image.rootFS?.type"
            >{{ image.rootFS?.type }}</x-description-item>
            <x-description-item
              :label="t('fields.comment')"
              v-if="image.comment"
              :span="2"
            >{{ image.comment }}</x-description-item>
          </x-description>
          <x-panel :title="t('fields.layers')" v-if="image.histories && image.histories.length">
            <n-data-table
              remote
              size="small"
              :columns="columns"
              :data="image.histories"
              scroll-x="max-content"
            />
          </x-panel>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="raw" :tab="t('fields.raw')">
        <x-code :code="raw" language="json" />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
  NDataTable,
  NTabs,
  NTabPane,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XCode from "@/components/Code.vue";
import XPanel from "@/components/Panel.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import imageApi from "@/api/image";
import type { Image } from "@/api/image";
import { useRoute } from "vue-router";
import { formatSize, renderTags } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const image = ref({} as Image);
const raw = ref('');
const node = route.params.node as string || '';
const columns = [
  {
    title: t('fields.sn'),
    key: "no",
    width: 45,
    fixed: "left" as const,
    render: (h: any, i: number) => i + 1,
  },
  {
    title: t('fields.instruction'),
    key: "createdBy",
    width: 500,
  },
  {
    title: t('fields.tags'),
    key: "image",
    render(i: Image) {
      if (i.tags) {
        return renderTags(i.tags?.map(t => {
          return { text: t, type: 'default' }
        }), true, 6)
      }
    },
  },
  {
    title: t('fields.size'),
    key: "size",
    width: 90,
    render(i: Image) {
      return formatSize(i.size)
    }
  },
  {
    title: t('fields.comment'),
    key: "comment",
  },
  {
    title: t('fields.created_at'),
    key: "createdAt",
    width: 150,
  },
];

async function fetchData() {
  const id = route.params.id as string;
  let r = await imageApi.find(node, id);
  raw.value = r.data?.raw as string;
  image.value = r.data?.image as Image;
  image.value.histories && image.value.histories.reverse();
}

onMounted(fetchData);
</script>