import { api } from './client';
import type { Family, User, FamilyWithOwner } from '../types';

export interface CreateFamilyRequest {
  name: string;
  owner_username: string;
  owner_password: string;
  owner_display_name: string;
}

export interface AddMemberRequest {
  username: string;
  pin: string;
  display_name: string;
}

export const familiesApi = {
  // Super admin
  listAll: () => api.get<FamilyWithOwner[]>('/admin/families'),
  create: (data: CreateFamilyRequest) => api.post<FamilyWithOwner>('/admin/families', data),
  deleteFamily: (id: number) => api.delete<void>(`/admin/families/${id}`),

  // Family owner/member
  getOwn: () => api.get<Family>('/family/'),
  listMembers: () => api.get<User[]>('/family/members'),
  addMember: (data: AddMemberRequest) => api.post<User>('/family/members', data),
  updateMember: (id: number, data: Partial<AddMemberRequest & { is_active: boolean }>) =>
    api.put<User>(`/family/members/${id}`, data),
  deactivateMember: (id: number) => api.delete<void>(`/family/members/${id}`),
  reactivateMember: (id: number) => api.post<void>(`/family/members/${id}/reactivate`, {}),
};
