<section class="min-h-screen grid grid-cols-1 md:grid-cols-2">
  <!-- Left: Login Form -->
  <div class="flex items-center justify-center px-6 py-12 bg-white">
    <div class="w-full max-w-sm">
      <div class="mb-6 text-center">
        <h1 class="text-3xl font-semibold text-gray-900">Sign in to your account</h1>
        <p class="text-sm text-gray-600 mt-2">Welcome back. Please enter your details.</p>
      </div>

      <form class="space-y-5" method="POST">
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            id="email"
            name="email"
            type="email"
            required
            placeholder="you@example.com"
            class="w-full px-4 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-black focus:border-black"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            required
            placeholder="••••••••"
            class="w-full px-4 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-black focus:border-black"
          />
        </div>

        <div class="flex items-center justify-between">
          <label class="flex items-center space-x-2 text-sm text-gray-600">
            <input type="checkbox" class="accent-black" />
            <span>Remember me</span>
          </label>
          <a href="/forgot" class="text-sm text-gray-500 hover:text-black">Forgot password?</a>
        </div>

        <button
          type="submit"
          class="w-full bg-black text-white py-2 rounded-md text-sm font-medium hover:bg-gray-800 transition-colors"
        >
          Sign In
        </button>
      </form>

      <p class="text-center text-sm text-gray-500 mt-6">
        Don't have an account?
        <a href="/register" class="text-black font-medium hover:underline">Sign up</a>
      </p>
    </div>
  </div>

  <!-- Right: Image / Visual Placeholder -->
  <div class="hidden md:block bg-gray-100">
    <!-- Replace this div with <img src="/your-image.jpg" alt="" class="w-full h-full object-cover" /> -->
    <div class="w-full h-full flex items-center justify-center text-gray-400 text-xl">
      Image Placeholder
    </div>
  </div>
</section>
