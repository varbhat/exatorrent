<script lang="ts">
  // Password box visibility toggle
  let pwvisible = false;
  let pwbox: HTMLInputElement;
  function toggleinput() {
    pwvisible = !pwvisible;
    pwbox.type = pwvisible ? 'text' : 'password';
  }
  // Signing in
  let exausername: string;
  let exapassword: string;
  import { Connect } from './core';
  import { onMount } from 'svelte';
  onMount(() => {
    let un = localStorage.getItem('exausername');
    let pw = localStorage.getItem('exapassword');
    if (un != '' && un != undefined && un != null) {
      if (pw != '' && pw != undefined && pw != null) {
        Connect();
      } else {
        return;
      }
    } else {
      return;
    }
  });
  function signIn() {
    if (exausername != '' && exausername != undefined && exausername != null) {
      if (exapassword != '' && exapassword != undefined && exapassword != null) {
        if (!(exausername.length > 5) || !(exapassword.length > 5)) {
          alert('Invalid Credentials');
          return;
        }
        localStorage.setItem('exausername', exausername);
        localStorage.setItem('exapassword', exapassword);
        Connect();
      } else {
        alert('Password Field Cannot be Empty');
        return;
      }
    } else {
      alert('Username Field Cannot be Empty');
      return;
    }
    return;
  }

  let entertosignin = (event: KeyboardEvent) => {
    if (event.code === 'Enter') {
      signIn();
    }
  };
</script>

<div class="mt-10  flex items-center justify-center px-4">
  <div class="max-w-md w-full ">
    <div>
      <a href="https://github.com/varbhat/exatorrent" target="_blank" rel="noopener noreferrer">
        <h2 class=" text-center text-5xl font-extrabold text-blue-200">exatorrent</h2>
      </a>
      <p class="mt-2 text-center text-sm text-neutral-300">Sign in to your account</p>
    </div>

    <div class="mt-2">
      <label for="username" class="sr-only">Username</label>

      <input id="username" name="email" type="text" bind:value={exausername} required class="bg-neutral-800 my-2 appearance-none rounded-md w-full px-3 py-2 border border-neutral-800 placeholder-neutral-500 text-neutral-200 focus:outline-none" placeholder="Username" />

      <label for="password" class="sr-only">Password</label>

      <div class="flex bg-neutral-800 rounded-md my-2 appearance-none border border-neutral-800 w-full">
        <input
          id="password"
          name="password"
          type="password"
          bind:value={exapassword}
          bind:this={pwbox}
          on:keydown={entertosignin}
          required
          class=" bg-neutral-800 appearance-none rounded-md w-full flex-grow px-3 py-2  border-none placeholder-neutral-500 text-neutral-200  focus:outline-none"
          placeholder="Password" />
        <button type="button" class="focus:outline-none focus:text-green-500" on:click={toggleinput}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neutral-400  my-2 mx-2 flex-grow " fill="none" viewBox="0 0 24 24" stroke="currentColor">
            {#if pwvisible}
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
            {:else}
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            {/if}
          </svg>
        </button>
      </div>

      <button type="button" on:click={signIn} class="w-full my-2  py-2 px-4 border-none text-sm font-medium rounded-md text-white bg-blue-900  outline-none focus:outline-none"> Sign in </button>
    </div>
  </div>
</div>
