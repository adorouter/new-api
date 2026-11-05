/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { z } from 'zod'

import type { Model } from '../types'
import { parseModelTags as parseTagsFromUtils } from './model-utils'

// ============================================================================
// Model Form Schema
// ============================================================================

/**
 * Model form validation schema
 */
export const modelFormSchema = z.object({
  id: z.number().optional(),
  model_name: z.string().min(1, 'Model name is required'),
  description: z.string().default(''),
  icon: z.string().default(''),
  tags: z.array(z.string()).default([]),
  vendor_id: z.number(),
  endpoints: z.string().default(''),
  name_rule: z.number().min(0).max(3).default(0),
  status: z.boolean().default(true),
  sync_official: z.boolean().default(true),
  enable_groups: z.array(z.string()).default([]),
  quota_types: z.array(z.number()).default([]),
  // AdoRouter: the stored metadata JSON, carried through unchanged apart from
  // the two sizing fields below, so saving the form never drops other hints.
  metadata: z.string().default(''),
  context_window: z
    .string()
    .regex(/^\d*$/, 'Enter a whole number of tokens')
    .default(''),
  max_output_tokens: z
    .string()
    .regex(/^\d*$/, 'Enter a whole number of tokens')
    .default(''),
})

export type ModelFormValues = z.infer<typeof modelFormSchema>

export const MODEL_FORM_DEFAULT_VALUES: ModelFormValues = {
  model_name: '',
  description: '',
  icon: '',
  tags: [],
  vendor_id: 0,
  endpoints: '',
  name_rule: 0,
  status: true,
  sync_official: true,
  enable_groups: [],
  quota_types: [],
  metadata: '',
  context_window: '',
  max_output_tokens: '',
}

type MetadataObject = Record<string, unknown>

function parseMetadataObject(metadata: string): MetadataObject | null {
  if (!metadata.trim()) return {}
  try {
    const parsed: unknown = JSON.parse(metadata)
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      return parsed as MetadataObject
    }
  } catch {
    // Unreadable metadata is left exactly as stored.
  }
  return null
}

function sizingText(value: unknown): string {
  return typeof value === 'number' && Number.isInteger(value) && value > 0
    ? String(value)
    : ''
}

/**
 * The context window and output cap stored in a model's metadata JSON, as form text.
 */
export function readModelSizing(metadata: string): {
  context_window: string
  max_output_tokens: string
} {
  const parsed = parseMetadataObject(metadata) ?? {}
  return {
    context_window: sizingText(parsed.contextWindow),
    max_output_tokens: sizingText(parsed.maxOutputTokens),
  }
}

/**
 * Writes the two sizing fields into the metadata JSON and keeps every other key.
 * An empty field removes its key. Metadata that is not a JSON object is returned
 * unchanged rather than overwritten.
 */
export function writeModelSizing(
  metadata: string,
  contextWindow: string,
  maxOutputTokens: string
): string {
  const parsed = parseMetadataObject(metadata)
  if (parsed === null) return metadata
  const next: MetadataObject = { ...parsed }
  for (const [key, value] of [
    ['contextWindow', contextWindow],
    ['maxOutputTokens', maxOutputTokens],
  ] as const) {
    const tokens = Number(value)
    if (value !== '' && Number.isInteger(tokens) && tokens > 0) {
      next[key] = tokens
    } else {
      delete next[key]
    }
  }
  return Object.keys(next).length ? JSON.stringify(next) : ''
}

// ============================================================================
// Vendor Form Schema
// ============================================================================

/**
 * Vendor form validation schema
 */
export const vendorFormSchema = z.object({
  id: z.number().optional(),
  name: z.string().min(1, 'Vendor name is required'),
  description: z.string().default(''),
  icon: z.string().default(''),
  status: z.number().default(1),
})

export type VendorFormValues = z.infer<typeof vendorFormSchema>

// ============================================================================
// Form Data Transformation
// ============================================================================

/**
 * Transform model to form default values
 */
export function transformModelToFormDefaults(model: Model): ModelFormValues {
  return {
    id: model.id,
    model_name: model.model_name,
    description: model.description || '',
    icon: model.icon || '',
    tags: parseTagsFromUtils(model.tags),
    vendor_id: model.vendor_id ?? 0,
    endpoints: model.endpoints || '',
    name_rule: model.name_rule || 0,
    status: model.status === 1,
    sync_official: model.sync_official === 1,
    enable_groups: model.enable_groups || [],
    quota_types: model.quota_types || [],
    metadata: model.metadata || '',
    ...readModelSizing(model.metadata || ''),
  }
}

/**
 * Transform form data to model create/update payload
 */
export function transformFormDataToModelPayload(
  formData: ModelFormValues
): Partial<Model> {
  return {
    id: formData.id,
    model_name: formData.model_name,
    description: formData.description || '',
    icon: formData.icon || '',
    tags: formatTagsArray(formData.tags),
    vendor_id: formData.vendor_id,
    endpoints: formData.endpoints || '',
    name_rule: formData.name_rule,
    status: formData.status ? 1 : 0,
    sync_official: formData.sync_official ? 1 : 0,
    enable_groups: formData.enable_groups,
    quota_types: formData.quota_types,
    metadata: writeModelSizing(
      formData.metadata,
      formData.context_window,
      formData.max_output_tokens
    ),
  }
}

// ============================================================================
// Parsing and Formatting Helpers
// ============================================================================

/**
 * Format tags array to string
 */
export function formatTagsArray(tags: string[]): string {
  return tags.filter(Boolean).join(',')
}

/**
 * Validate JSON string
 */
export function validateJSON(value: string): boolean {
  if (!value || value.trim() === '') return true

  try {
    JSON.parse(value)
    return true
  } catch {
    return false
  }
}

/**
 * Validate endpoints JSON
 */
export function validateEndpoints(endpoints: string): boolean {
  return validateJSON(endpoints)
}
