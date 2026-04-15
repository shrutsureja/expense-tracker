import { api } from './client';
import type { Tag } from '../types';

export const tagsApi = {
  list: () => api.get<Tag[]>('/tags'),
  suggested: () => api.get<Tag[]>('/tags/suggested'),
  create: (name: string, icon: string) => api.post<Tag>('/tags', { name, icon }),
  delete: (id: number) => api.delete<void>(`/tags/${id}`),
};
