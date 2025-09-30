import { z } from "zod";

export const addAliasSchema = z.object({
  alias: z.string().min(1, "Alias is required"),
});

export type AddAliasSchema = z.infer<typeof addAliasSchema>;

export const websiteConfigSchema = z.object({
  websiteAccess: z.boolean(),
  indexDocument: z.string().nullish(),
  errorDocument: z.string().nullish(),
});

export type WebsiteConfigSchema = z.infer<typeof websiteConfigSchema>;

export const quotaSchema = z.object({
  enabled: z.boolean(),
  maxObjects: z.coerce.number().nullish(),
  maxSize: z.coerce.number().nullish(),
});

export type QuotaSchema = z.infer<typeof quotaSchema>;

export const allowKeysSchema = z.object({
  keys: z
    .object({
      checked: z.boolean(),
      keyId: z.string(),
      name: z.string(),
      read: z.boolean(),
      write: z.boolean(),
      owner: z.boolean(),
    })
    .array(),
});

export type AllowKeysSchema = z.infer<typeof allowKeysSchema>;
