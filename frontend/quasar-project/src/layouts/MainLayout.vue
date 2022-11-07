<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="bg-black">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title>
          Bookmark Management Tool
        </q-toolbar-title>

        <div>Index v.{{ 0.1 }}</div>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
    >
      <q-list>
        <q-item-label
          header
          class="MenuTitle"
        >
          Menu
        </q-item-label>
        <div class="menu-hr-line"></div>
        <EssentialLink
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
        />
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
import { defineComponent, ref } from 'vue'
import EssentialLink from 'components/EssentialLink.vue'

const linksList = [
  {
    title: 'Bookmark',
    caption: 'bookmark management',
    icon: 'bookmark',
    link: '/'
  },
  {
    title: 'Dashboard',
    caption: 'dashboard',
    icon: 'dashboard',
    link: ''
  },
  {
    title: 'Team',
    caption: 'Team Bookmark',
    icon: 'group',
    link: ''
  },
  {
    title: 'Accounts',
    caption: 'Manage Account',
    icon: 'manage_accounts',
    link: ''
  }
]

export default defineComponent({
  name: 'MainLayout',

  components: {
    EssentialLink
  },

  setup () {
    const leftDrawerOpen = ref(false)

    return {
      essentialLinks: linksList,
      leftDrawerOpen,
      toggleLeftDrawer () {
        leftDrawerOpen.value = !leftDrawerOpen.value
      }
    }
  }
})
</script>
<style scoped>
.MenuTitle{
  text-align: center;
  font-size: 18px;
  margin-top: 10px;
}
.menu-hr-line {
    border-top: 3px solid gainsboro;
    margin: 10px;
}
</style>
