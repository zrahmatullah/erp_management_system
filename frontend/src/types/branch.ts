export interface Branch {
  id: string;
  name: string;
  address: string;
  isActive: boolean;
}
export interface BranchSetting {
  branchId: string;
  settingKey: string;
  settingValue: string;
}
