<template>
  <div class="test-view">
    <div class="container">
      <h1>Test</h1>

      <div
        v-for="section in testSections"
        :key="section.path"
        class="test-section"
        @click="navigateTo(section.path)"
      >
        <h2>{{ section.title }}</h2>
        <p>{{ section.description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const testSections = [
  {
    title: 'Album Card Test',
    description: 'AlbumCardコンポーネントの表示テスト - サンプルアルバムとエッジケースの確認',
    path: '/test/album-card',
  },
  {
    title: 'Image Card Test',
    description: 'ImageCardコンポーネントの表示テスト - 画像の読み込みと表示の確認',
    path: '/test/image-card',
  },
  {
    title: 'traQ File Proxy Test',
    description: 'traQ File Proxy Testのテスト - ファイルの読み込みと表示の確認',
    path: '/test/traq-file',
  },
]

const navigateTo = (path: string) => {
  router.push(path)
}
</script>

<style scoped>
.test-view {
  min-height: 100vh;
  background: #f8f9fa;
  padding: 20px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

h1 {
  text-align: center;
  color: #333;
  margin-bottom: 40px;
  font-size: 32px;
}

.test-section {
  margin-bottom: 30px;
  padding: 20px;
  border-radius: 8px;
  cursor: pointer;
  transition:
    background-color 0.2s,
    box-shadow 0.2s;
}

.test-section:hover {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.test-section h2 {
  color: #555;
  margin-bottom: 10px;
  font-size: 24px;
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
}

.test-section p {
  color: #666;
  margin: 0;
  font-size: 14px;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .container {
    padding: 0 10px;
  }

  h1 {
    font-size: 28px;
  }

  .test-section h2 {
    font-size: 20px;
  }
}
</style>
