import './index.css';
import Index from './Index.svelte';
import { mount } from 'svelte';

const app = mount(Index, { target: document.getElementById('app')! });

export default app;
