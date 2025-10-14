import Link from 'next/link'

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Hero Section */}
      <section className="text-center py-20">
        <h1 className="text-5xl font-bold mb-4">
          Kube EC へようこそ
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Kubernetes + Go + Next.js で構築したモダンなECサイト
        </p>
        <Link
          href="/products"
          className="bg-primary text-white px-8 py-3 rounded-lg text-lg font-semibold hover:bg-blue-600 transition"
        >
          商品を見る
        </Link>
      </section>

      {/* Features Section */}
      <section className="grid md:grid-cols-3 gap-8 py-16">
        <div className="text-center p-6 rounded-lg border">
          <div className="text-4xl mb-4">🚀</div>
          <h3 className="text-xl font-semibold mb-2">高速配送</h3>
          <p className="text-gray-600">最短翌日お届け</p>
        </div>
        <div className="text-center p-6 rounded-lg border">
          <div className="text-4xl mb-4">🔒</div>
          <h3 className="text-xl font-semibold mb-2">安全な決済</h3>
          <p className="text-gray-600">複数の決済方法に対応</p>
        </div>
        <div className="text-center p-6 rounded-lg border">
          <div className="text-4xl mb-4">💯</div>
          <h3 className="text-xl font-semibold mb-2">品質保証</h3>
          <p className="text-gray-600">厳選された商品のみ</p>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-primary text-white rounded-lg p-12 text-center">
        <h2 className="text-3xl font-bold mb-4">今すぐ始めよう</h2>
        <p className="text-xl mb-6">アカウントを作成して特典をゲット</p>
        <Link
          href="/register"
          className="bg-white text-primary px-8 py-3 rounded-lg text-lg font-semibold hover:bg-gray-100 transition inline-block"
        >
          無料登録
        </Link>
      </section>
    </div>
  )
}
