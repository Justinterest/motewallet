import type { User } from "@/types/auth";

/** Merchant can use trading features (deposit, withdraw, etc.). */
export function canMerchantOperate(user: User | null | undefined): boolean {
  return user?.status === "ACTIVE";
}

export interface VerificationBannerInfo {
  message: string;
  ctaLabel: string;
  ctaHref: string;
  variant: "warning" | "info" | "error";
}

export function getVerificationBanner(
  user: User | null | undefined
): VerificationBannerInfo | null {
  if (!user || canMerchantOperate(user)) {
    return null;
  }

  switch (user.status) {
    case "PENDING_AGREEMENT":
      return {
        message: "请先签署平台服务协议，完成后即可提交企业实名认证。",
        ctaLabel: "去签署协议",
        ctaHref: "/onboarding/agreement",
        variant: "warning",
      };
    case "PENDING_KYC":
      return {
        message: "请完成企业实名认证。认证通过前可浏览各功能页面，但无法发起资金操作。",
        ctaLabel: "去实名认证",
        ctaHref: "/onboarding/kyc",
        variant: "warning",
      };
    case "KYC_REVIEWING":
      return {
        message: "企业实名认证审核中，请耐心等待。审核期间仅可查看，无法操作资金。",
        ctaLabel: "查看认证进度",
        ctaHref: "/onboarding/status",
        variant: "info",
      };
    case "KYC_FAILED":
      return {
        message: "企业实名认证未通过，请修改资料后重新提交。",
        ctaLabel: "重新认证",
        ctaHref: "/onboarding/kyc",
        variant: "error",
      };
    default:
      if (user.kyc_status === "AUTHING") {
        return {
          message: "企业实名认证审核中，请耐心等待。",
          ctaLabel: "查看认证进度",
          ctaHref: "/onboarding/status",
          variant: "info",
        };
      }
      return {
        message: "账户尚未完成认证，部分功能暂不可用。",
        ctaLabel: "完成认证",
        ctaHref: "/onboarding/kyc",
        variant: "warning",
      };
  }
}
