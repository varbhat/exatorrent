<script lang="ts">
  import { Send, userconnlist, userlist, isAdmin } from './core';
  import slocation from 'slocation';
  import { onMount } from 'svelte';
  import Useredit from './Useredit.svelte';

  onMount(() => {
    if ($isAdmin === false) {
      slocation.goto('/');
    }
  });

  let userconnlistOpen = false;
  let manageUsersOpen = false;
  let addUserOpen = false;
  let newusername = '';
  let newpassword = '';
  let newusertype: 'user' | 'disabled' | 'admin' = 'user';
  let pwbox: HTMLInputElement;
  let pwvisible = false;

  let toggleinput = () => {
    pwvisible = !pwvisible;
    pwbox.type = pwvisible ? 'text' : 'password';
  };

  let addUser = () => {
    Send({
      command: 'adduser',
      data1: newusername,
      data2: newpassword,
      data3: newusertype,
      aop: 1
    });
  };

  let typenotostring = (t: number): 'user' | 'disabled' | 'admin' => {
    switch (t) {
      case 0:
        return 'user';
      case 1:
        return 'admin';
      case -1:
        return 'disabled';
    }
  };

  let userconnlistaction = () => {
    if (userconnlistOpen === false) {
      Send({
        command: 'listuserconns',
        aop: 1
      });
      userconnlistOpen = true;
    } else {
      userconnlistOpen = false;
    }
  };

  let manageUsersaction = () => {
    if (manageUsersOpen === false) {
      Send({
        command: 'getusers',
        aop: 1
      });
      manageUsersOpen = true;
    } else {
      manageUsersOpen = false;
    }
  };

  let useraddaction = () => {
    if (addUserOpen === false) {
      addUserOpen = true;
    } else {
      addUserOpen = false;
    }
  };
</script>

<div class="mx-auto max-w-3xl ">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={userconnlistaction}>
        <p class="ml-3 font-medium  truncate">User Connections</p>
      </div>

      {#if userconnlistOpen === true}
        <button
          type="button"
          class="flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0"
          on:click={() => {
            Send({
              command: 'listuserconns',
              aop: 1
            });
          }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      {/if}

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1" on:click={userconnlistaction}>
        {#if userconnlistOpen === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        {/if}
      </button>
    </div>

    <div class="flex flex-col">
      {#if userconnlistOpen === true}
        {#each $userconnlist as eachuser (eachuser?.username)}
          <div class="text-neutral-200 bg-neutral-900 px-1 py-3 rounded-md w-full my-1">
            <div>
              <div class="flex items-center justify-between flex-wrap py-1 mr-1">
                <div class="w-0 flex-1 flex">
                  <p class="font-medium  break-all mx-1">
                    {eachuser?.username}
                    {#if eachuser?.isadmin === true}(admin){/if}
                  </p>
                </div>

                <div
                  class="-mr-1 flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1"
                  on:click={() => {
                    Send({
                      command: 'kickuser',
                      data1: eachuser?.username,
                      aop: 1
                    });
                  }}>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
              </div>
            </div>
            <p class="text-sm  font-extralight truncate mx-1">
              Connected at {new Date(eachuser?.contime)?.toLocaleString()}
            </p>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>

<div class="mx-auto max-w-3xl ">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={manageUsersaction}>
        <p class="ml-3 font-medium  truncate">Manage Users</p>
      </div>

      {#if manageUsersOpen === true}
        <button
          type="button"
          class="flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0"
          on:click={() => {
            Send({
              command: 'getusers',
              aop: 1
            });
          }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      {/if}

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1" on:click={manageUsersaction}>
        {#if manageUsersOpen === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        {/if}
      </button>
    </div>

    <div class="flex flex-col">
      {#if manageUsersOpen === true}
        {#each $userlist as eachuser (eachuser?.Username)}
          <Useredit username={eachuser?.Username} selected={typenotostring(eachuser?.UserType)} timestring={new Date(eachuser?.CreatedAt)?.toLocaleString()} />
        {/each}
      {/if}
    </div>
  </div>
</div>

<div class="mx-auto max-w-3xl ">
  <div class="bg-black grid grid-flow-row text-white rounded-lg m-3 p-2 cursor-pointer focus:outline-none focus-within:bg-black noHL">
    <div class="flex items-center justify-between flex-wrap py-1 px-3">
      <div class="w-0 flex-1 flex items-center" on:click={useraddaction}>
        <p class="ml-3 font-medium  truncate">Add User</p>
      </div>

      <button type="button" class="-mr-1 flex p-2 rounded-md bg-neutral-800 focus:outline-none flex-shrink-0 mx-1" on:click={useraddaction}>
        {#if addUserOpen === true}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        {/if}
      </button>
    </div>

    <div class="flex flex-col">
      {#if addUserOpen === true}
        <div class="mt-2">
          <label for="username" class="sr-only">Username</label>

          <input id="username" name="email" type="text" bind:value={newusername} required class="bg-neutral-800 appearance-none rounded-md w-full px-3 py-2 border border-neutral-800 placeholder-neutral-500 text-neutral-200 focus:outline-none" placeholder="Username" />

          <label for="password" class="sr-only">Password</label>

          <div class="flex bg-neutral-800 rounded-md my-2 appearance-none border border-neutral-800 w-full">
            <input
              id="password"
              name="password"
              type="password"
              autocomplete="current-password"
              bind:value={newpassword}
              bind:this={pwbox}
              required
              class="bg-neutral-800 appearance-none rounded-md w-full flex-grow px-3 py-2 border-none placeholder-neutral-500 text-neutral-200 focus:outline-none"
              placeholder="Password" />
            <button type="button" class="focus:outline-none focus:text-green-500" on:click={toggleinput}>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neutral-400 my-2 mx-2 flex-grow" fill="none" viewBox="0 0 24 24" stroke="currentColor">
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

          <select class="bg-neutral-800 rounded-md w-full flex-grow px-3 py-2 border-none placeholder-neutral-500 text-neutral-200 focus:outline-none" bind:value={newusertype}>
            <option value="user">User</option>
            <option value="disabled">Disabled</option>
            <option value="admin">Admin</option>
          </select>

          <button type="button" on:click={addUser} class="w-full my-2 py-2 px-4 border-none text-sm font-medium rounded-md text-white bg-indigo-900 outline-none focus:outline-none">Add User</button>
        </div>
      {/if}
    </div>
  </div>
</div>
