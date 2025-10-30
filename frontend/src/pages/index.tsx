import Head from 'next/head'

export default function Home() {
  return (
    <>
      <Head>
        <title>Task Management App</title>
        <meta name="description" content="Task Management Application" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </Head>
      <main className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto py-12 px-4">
          <div className="text-center">
            <h1 className="text-4xl font-bold text-gray-900 mb-4">
              Task Management App
            </h1>
            <p className="text-lg text-gray-600 mb-8">
              Frontend coming in Week 6
            </p>
          </div>
        </div>
      </main>
    </>
  )
}
