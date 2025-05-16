import type { Actions } from './$types';

export const actions = {
      default: async ({ request }) => {
            const formData = await request.formData();
            const values = formData.getAll("options");

            console.log(values)
            return { success: true }
      }
} satisfies Actions;