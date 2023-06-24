<script lang="ts">
  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import { goto } from "$app/navigation";
  import { browser } from "$app/environment";
  import { PUBLIC_BACKEND_API_BASE } from "$env/static/public";

  const templates = writable<{
    id: number;
    name: string;
  }[]>();

  onMount(async () => {
    if (!browser) {
      return;
    }

    fetch(`${PUBLIC_BACKEND_API_BASE}/api/templates`)
      .then(response => response.json())
      .then(data => templates.set(data))
      .catch(() => []);
  });
</script>

<div class="flex flex-row gap-4 justify-center mb-8 w-full max-w-screen-2xl">
  <input class="input p-2 w-1/2" type="text" placeholder="Search..." />
  <button type="button" class="btn variant-ghost" on:click={() => goto("/new-template")}>New Template</button>
</div>

<div class="grid grid-cols-5 gap-3 w-full max-w-screen-2xl 2xl:grid-cols-5">
  {#if $templates}
    {#each $templates as template}
      <a class="block card variant-soft card-hover w-full" href="/template/{ template.id }">
        <header>
          <div class="w-full h-64 bg-contain bg-no-repeat bg-center" style="background-image: url('{ PUBLIC_BACKEND_API_BASE }/template/{ template.id }.webp')"></div>
        </header>
        <section class="p-3"><span class="h3">{ template.name }</span></section>
      </a>
    {/each}
  {/if}
</div>
