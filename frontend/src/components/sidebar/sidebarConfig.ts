export type SidebarItem = {
  title: string
  icon?: string
  path: string
}

export const sidebarItems: SidebarItem[] = [
  { title: 'Home', icon: '/dummyHomeIcon.png', path: '/' },
  { title: 'Albums', icon: '/dummyAlbumsIcon.png', path: '/albums' },
  { title: 'Saved', icon: '/dummySavedIcon.png', path: '/saved' },
  { title: 'Account', path: '/account' },
]
