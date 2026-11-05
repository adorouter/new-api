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
import { useFormContext } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

import type { ChannelFormValues } from '../lib/channel-form'

// AdoRouter: the provider facts behind the public site's provider labels, the
// "may train" notice and the consent dialog. They live in the channel's
// setting.data_policy, which the console could not edit before.
const PROMPTS_OPTIONS = [
  { value: 'unknown', label: 'Not published' },
  { value: 'none', label: 'Neither stores nor trains' },
  { value: 'retains', label: 'Stores prompts' },
  { value: 'trains', label: 'Trains on prompts' },
] as const

const TRANSFER_OPTIONS = [
  { value: 'unknown', label: 'Unknown' },
  { value: 'adequacy', label: 'EU adequacy decision' },
  { value: 'dpf', label: 'EU-US Data Privacy Framework' },
  { value: 'scc', label: 'Standard contractual clauses' },
  { value: 'none', label: 'No safeguard' },
] as const

export function DataPolicySetting(props: { disabled?: boolean }) {
  const { t } = useTranslation()
  const form = useFormContext<ChannelFormValues>()

  return (
    <div className='space-y-4'>
      <div className='space-y-1'>
        <p className='text-sm font-medium'>{t('Provider data policy')}</p>
        <p className='text-muted-foreground text-sm'>
          {t(
            'Who runs this upstream and what it does with prompts. The public model pages and the consent dialog read these facts.'
          )}
        </p>
      </div>
      <FormField
        control={form.control}
        name='data_policy_provider'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Provider company')}</FormLabel>
            <FormControl>
              <Input
                placeholder={t('e.g., Example AI Ltd')}
                disabled={props.disabled}
                {...field}
                value={field.value ?? ''}
              />
            </FormControl>
            <FormDescription>
              {t(
                'Channels with the same company share one consent. Never shown publicly.'
              )}
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name='data_policy_seat_country'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Provider country')}</FormLabel>
            <FormControl>
              <Input
                placeholder='US'
                maxLength={2}
                autoCapitalize='characters'
                disabled={props.disabled}
                {...field}
                value={field.value ?? ''}
                onChange={(event) =>
                  field.onChange(event.target.value.toUpperCase())
                }
              />
            </FormControl>
            <FormDescription>
              {t('Two-letter country code where the provider is based.')}
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name='data_policy_prompts'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('What it does with prompts')}</FormLabel>
            <Select
              disabled={props.disabled}
              items={PROMPTS_OPTIONS.map((option) => ({
                value: option.value,
                label: t(option.label),
              }))}
              value={field.value || 'unknown'}
              onValueChange={(value) => field.onChange(value ?? 'unknown')}
            >
              <FormControl>
                <SelectTrigger aria-label={t('What it does with prompts')}>
                  <SelectValue />
                </SelectTrigger>
              </FormControl>
              <SelectContent alignItemWithTrigger={false}>
                <SelectGroup>
                  {PROMPTS_OPTIONS.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {t(option.label)}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormDescription>
              {t(
                'Choose "Neither stores nor trains" only when the provider publishes that policy.'
              )}
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name='data_policy_transfer'
        render={({ field }) => (
          <FormItem>
            <FormLabel>{t('Transfer safeguard')}</FormLabel>
            <Select
              disabled={props.disabled}
              items={TRANSFER_OPTIONS.map((option) => ({
                value: option.value,
                label: t(option.label),
              }))}
              value={field.value || 'unknown'}
              onValueChange={(value) => field.onChange(value ?? 'unknown')}
            >
              <FormControl>
                <SelectTrigger aria-label={t('Transfer safeguard')}>
                  <SelectValue />
                </SelectTrigger>
              </FormControl>
              <SelectContent alignItemWithTrigger={false}>
                <SelectGroup>
                  {TRANSFER_OPTIONS.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {t(option.label)}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormDescription>
              {t(
                'The legal basis for sending EU users’ data to the provider’s country.'
              )}
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
    </div>
  )
}
