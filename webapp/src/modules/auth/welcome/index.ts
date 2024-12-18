import { createSessionStorageHandlers } from '@/utils/sessionStorage';
import WelcomeBanner from './welcomeBanner.svelte';

const WelcomeSession = createSessionStorageHandlers('welcome');

export { WelcomeBanner, WelcomeSession };
