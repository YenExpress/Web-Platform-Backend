import crypto from 'crypto'
import { v4 as uuid } from 'uuid';

export async function createUniqueToken(baseByteSize: number) {
    const tokenBytes = crypto.randomBytes(baseByteSize)
    const token = tokenBytes.toString('hex') + uuid()
    return token
}

export function calculateExpiryDateFromString(expirationString: string): string {
	const match = expirationString.match(/^(\d+)\s+(d|m|h)$/i);
	if (!match) {
	  throw new Error('Invalid input format');
	}
  
	const expiration = parseInt(match[1]);
	const unit = match[2].toLowerCase();
  
	const now = new Date();
	switch (unit) {
	  case 'm':
		now.setMinutes(now.getMinutes() + expiration);
		break;
	  case 'h':
		now.setHours(now.getHours() + expiration);
		break;
	  case 'd':
		now.setDate(now.getDate() + expiration);
		break;
	  default:
		throw new Error('Invalid unit');
	}
	return now.toISOString();
  }