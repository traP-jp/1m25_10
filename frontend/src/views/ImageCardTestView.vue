<template>
  <div class="image-card-test">
    <div class="container">
      <h1>ImageCard Component Test</h1>

      <div class="test-section">
        <h2>Sample Images</h2>
        <div class="image-grid">
          <ImageCard
            v-for="image in sampleImages"
            :key="image.id"
            :url="image.url"
            :alt="image.alt"
          />
        </div>
      </div>

      <div class="test-section">
        <h2>Different Image Sizes & Types</h2>
        <div class="image-grid">
          <ImageCard
            v-for="image in variousImages"
            :key="image.id"
            :url="image.url"
            :alt="image.alt"
          />
        </div>
      </div>

      <div class="test-section">
        <h2>Error Handling (Broken URLs)</h2>
        <div class="image-grid">
          <ImageCard
            v-for="image in brokenImages"
            :key="image.id"
            :url="image.url"
            :alt="image.alt"
          />
        </div>
      </div>

      <div class="test-section">
        <h2>Usage Instructions</h2>
        <div class="instructions">
          <h3>ImageCard Features:</h3>
          <ul>
            <li>✅ Click to select/deselect images</li>
            <li>⌨️ Press Enter or Space to toggle selection</li>
            <li>🖱️ Hover effects and smooth transitions</li>
            <li>📱 Responsive design (196px × 196px)</li>
            <li>🔄 Automatic fallback for broken images</li>
            <li>🎨 Consistent design with project theme</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ImageCard from '@/components/ImageCard.vue'

interface TestImage {
  id: string
  url: string
  alt: string
}

// サンプル画像データ（実際の画像URL）
const sampleImages: TestImage[] = [
  {
    id: '1',
    url: 'https://picsum.photos/400/400?random=1',
    alt: 'ランダム画像1',
  },
  {
    id: '2',
    url: 'https://picsum.photos/400/400?random=2',
    alt: 'ランダム画像2',
  },
  {
    id: '3',
    url: 'https://picsum.photos/400/400?random=3',
    alt: 'ランダム画像3',
  },
  {
    id: '4',
    url: 'https://picsum.photos/400/400?random=4',
    alt: 'ランダム画像4',
  },
  {
    id: '5',
    url: 'https://picsum.photos/400/400?random=5',
    alt: 'ランダム画像5',
  },
  {
    id: '6',
    url: 'https://picsum.photos/400/400?random=6',
    alt: 'ランダム画像6',
  },
]

// 様々なサイズ・タイプの画像
const variousImages: TestImage[] = [
  {
    id: 'landscape',
    url: 'https://picsum.photos/800/400?random=10',
    alt: '横長画像',
  },
  {
    id: 'portrait',
    url: 'https://picsum.photos/400/800?random=11',
    alt: '縦長画像',
  },
  {
    id: 'square',
    url: 'https://picsum.photos/500/500?random=12',
    alt: '正方形画像',
  },
  {
    id: 'small',
    url: 'https://picsum.photos/200/200?random=13',
    alt: '小さい画像',
  },
  {
    id: 'large',
    url: 'https://picsum.photos/1200/1200?random=14',
    alt: '大きい画像',
  },
  {
    id: 'grayscale',
    url: 'https://picsum.photos/400/400?random=15&grayscale',
    alt: 'グレースケール画像',
  },
]

// 壊れたURL（エラーハンドリングテスト用）
const brokenImages: TestImage[] = [
  {
    id: 'broken1',
    url: 'https://example.com/nonexistent-image.jpg',
    alt: '存在しない画像1',
  },
  {
    id: 'broken2',
    url: 'https://invalid-url/image.png',
    alt: '無効なURL',
  },
  {
    id: 'broken3',
    url: '/nonexistent-local-image.jpg',
    alt: 'ローカルの存在しない画像',
  },
  {
    id: 'broken4',
    url: '',
    alt: '空のURL',
  },
]
</script>

<style scoped>
.image-card-test {
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
  margin-bottom: 50px;
}

.test-section h2 {
  color: #555;
  margin-bottom: 20px;
  font-size: 24px;
  border-bottom: 2px solid #005bac;
  padding-bottom: 10px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, 196px);
  gap: 20px;
  padding: 20px 0;
  justify-content: start;
}

.instructions {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  max-width: 600px;
}

.instructions h3 {
  color: #005bac;
  margin-bottom: 16px;
  font-size: 18px;
}

.instructions ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.instructions li {
  padding: 8px 0;
  border-bottom: 1px solid #eee;
  font-size: 14px;
  line-height: 1.4;
}

.instructions li:last-child {
  border-bottom: none;
}

@media (max-width: 768px) {
  .image-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 15px;
  }

  .container {
    padding: 0 10px;
  }

  h1 {
    font-size: 24px;
  }

  .test-section h2 {
    font-size: 20px;
  }
}

@media (max-width: 480px) {
  .image-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 10px;
  }
}
</style>
