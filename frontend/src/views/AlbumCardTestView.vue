<template>
  <div class="album-card-test">
    <div class="container">
      <h1>AlbumCard Component Test</h1>

      <div class="test-section">
        <h2>Sample Albums</h2>
        <div class="album-grid">
          <AlbumCard
            v-for="album in sampleAlbums"
            :key="album.id"
            :id="album.id"
            :title="album.title"
            :image-urls="getImageUrls(album.images)"
            :created-at="album.created_at"
            :creator="getCreatorData(album.creator)"
          />
        </div>
      </div>

      <div class="test-section">
        <h2>Edge Cases</h2>
        <div class="album-grid">
          <AlbumCard
            :id="emptyAlbum.id"
            :title="emptyAlbum.title"
            :image-urls="getImageUrls(emptyAlbum.images)"
            :created-at="emptyAlbum.created_at"
            :creator="getCreatorData(emptyAlbum.creator)"
          />
          <AlbumCard
            :id="singleImageAlbum.id"
            :title="singleImageAlbum.title"
            :image-urls="getImageUrls(singleImageAlbum.images)"
            :created-at="singleImageAlbum.created_at"
            :creator="getCreatorData(singleImageAlbum.creator)"
          />
          <AlbumCard
            :id="manyImagesAlbum.id"
            :title="manyImagesAlbum.title"
            :image-urls="getImageUrls(manyImagesAlbum.images)"
            :created-at="manyImagesAlbum.created_at"
            :creator="getCreatorData(manyImagesAlbum.creator)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AlbumCard from '@/components/AlbumCard.vue'
import type { Album } from '@/types/album'

// クリエーターデータのマッピング
const creators = {
  user1: {
    id: 'quarantineeeeeeeeee',
    name: 'くあらんてぃん',
    avatarUrl: 'https://q.trap.jp/api/v3/public/icon/quarantineeeeeeeeee',
  },
  user2: {
    id: 'user2',
    name: 'Alice Johnson',
    avatarUrl: 'https://picsum.photos/100/100?random=2',
  },
  user3: { id: 'user3', name: '佐藤花子', avatarUrl: 'https://picsum.photos/100/100?random=3' },
  user4: { id: 'user4', name: 'Michael Brown', avatarUrl: undefined }, // アバターなしのテスト
  user5: { id: 'user5', name: '山田次郎', avatarUrl: 'https://picsum.photos/100/100?random=5' },
  user6: { id: 'user6', name: 'Sarah Wilson', avatarUrl: 'https://picsum.photos/100/100?random=6' },
  user7: { id: 'user7', name: '鈴木三郎', avatarUrl: 'https://picsum.photos/100/100?random=7' },
}

// クリエーターIDからクリエーターオブジェクトを取得する関数
const getCreatorData = (creatorId: string) => {
  return (
    creators[creatorId as keyof typeof creators] || {
      id: creatorId,
      name: 'Unknown User',
      avatarUrl: undefined,
    }
  )
}

// 画像IDから画像URLを生成する関数
const getImageUrls = (imageIds: string[]): string[] => {
  return imageIds.map((id) => {
    // 実際のAPIエンドポイントの代わりに、Lorem Picsumを使用
    // ランダムシードとしてidのハッシュ値を使用
    const seed = id.split('').reduce((a, b) => {
      a = (a << 5) - a + b.charCodeAt(0)
      return a & a
    }, 0)
    return `https://picsum.photos/400/400?random=${Math.abs(seed)}`
  })
}

// サンプルアルバムデータ
const sampleAlbums: Album[] = [
  {
    id: '1',
    title: '夏の思い出',
    description: '今年の夏に撮った写真たちです。海や山での素敵な瞬間を集めました。',
    creator: 'user1',
    images: ['img1', 'img2', 'img3', 'img4'],
    created_at: '2025-08-15T10:30:00Z',
    updated_at: '2025-08-15T10:30:00Z',
  },
  {
    id: '2',
    title: '旅行写真コレクション',
    description: '世界各地を旅した時の写真コレクションです。',
    creator: 'user2',
    images: ['img5', 'img6'],
    created_at: '2025-07-20T14:15:00Z',
    updated_at: '2025-07-20T14:15:00Z',
  },
  {
    id: '3',
    title: 'ペットの写真',
    description: 'かわいい愛犬の成長記録です。',
    creator: 'user3',
    images: ['img7', 'img8', 'img9'],
    created_at: '2025-06-10T09:00:00Z',
    updated_at: '2025-06-10T09:00:00Z',
  },
  {
    id: '4',
    title: 'とても長いタイトルのアルバムテストとても長いタイトルのアルバムテスト',
    description:
      'とても長い説明文のテストです。この説明文は非常に長く、複数行に渡って表示されることを想定しています。文字数制限のテストも兼ねています。',
    creator: 'user4',
    images: ['img10', 'img11', 'img12', 'img13', 'img14', 'img15'],
    created_at: '2025-05-01T16:45:00Z',
    updated_at: '2025-05-01T16:45:00Z',
  },
  {
    id: '5',
    title: 'フォトグラフィー作品集',
    description: 'プロの写真家による美しい風景写真の数々。',
    creator: 'user2',
    images: ['img16', 'img17'],
    created_at: '2025-07-10T11:20:00Z',
    updated_at: '2025-07-10T11:20:00Z',
  },
  {
    id: '6',
    title: '日常のスナップ',
    description: '何気ない日常の瞬間を切り取った写真たち。',
    creator: 'user1',
    images: ['img18', 'img19', 'img20', 'img21', 'img22'],
    created_at: '2025-08-05T15:45:00Z',
    updated_at: '2025-08-05T15:45:00Z',
  },
]

// エッジケースのテスト用アルバム
const emptyAlbum: Album = {
  id: 'empty',
  title: '空のアルバム',
  description: '画像がないアルバムです。',
  creator: 'user5',
  images: [],
  created_at: '2025-09-01T12:00:00Z',
  updated_at: '2025-09-01T12:00:00Z',
}

const singleImageAlbum: Album = {
  id: 'single',
  title: '1枚だけのアルバム',
  description: '画像が1枚だけのアルバムです。',
  creator: 'user6',
  images: ['single-img'],
  created_at: '2025-08-25T08:30:00Z',
  updated_at: '2025-08-25T08:30:00Z',
}

const manyImagesAlbum: Album = {
  id: 'many',
  title: 'たくさんの画像があるアルバム（10枚以上）',
  description: '10枚以上の画像があるアルバムです。',
  creator: 'user7',
  images: Array.from({ length: 12 }, (_, i) => `many-img-${i + 1}`),
  created_at: '2025-08-30T20:00:00Z',
  updated_at: '2025-08-30T20:00:00Z',
}
</script>

<style scoped>
.album-card-test {
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
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
}

.album-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  padding: 20px 0;
}

@media (max-width: 768px) {
  .album-grid {
    grid-template-columns: 1fr;
  }

  .container {
    padding: 0 10px;
  }
}
</style>
