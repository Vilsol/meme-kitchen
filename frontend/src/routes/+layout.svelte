<script>
  import '@skeletonlabs/skeleton/themes/theme-gold-nouveau.css';
  import '@skeletonlabs/skeleton/styles/skeleton.css';
  import '../app.postcss';

  import { Toast, toastStore } from '@skeletonlabs/skeleton';
  import copy from 'copy-to-clipboard';

  const copyUrl = () => {
    copy('https://meme.kitchen');
    toastStore.trigger({
      message: 'Copied to clipboard!'
    });
  }
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
  <link href="https://fonts.googleapis.com/css2?family=Amatic+SC&display=swap" rel="stylesheet">

  <meta property="og:title" content="meme.kitchen"/>
  <meta property="og:image" content="/banner.jpg">
  <meta property="og:description" content="Dynamic meme generator">
</svelte:head>

<Toast />

<div class="flex flex-col items-center h-full scroll-auto pt-16">
  <div class="flex mb-8">
    <button tabindex="0" on:click={copyUrl}>
      <h1 class="text-8xl cursor-pointer select-none hover:bg-slate-500/5 active:bg-slate-500/10 rounded p-3 title">
        <span class="text-neutral-400/50">//</span>meme.kitchen
      </h1>
    </button>
  </div>

  <slot />
</div>

<style lang="postcss">
  .title {
    font-family: 'Amatic SC', sans-serif;
  }
</style>