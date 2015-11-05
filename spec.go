package lowrider

// http://infinispan.org/docs/8.0.x/user_guide/user_guide.html#_hot_rod_protocol_2_4

/***********************************************************
 * Data Types
 *
 * All key and values are sent and stored as byte arrays.
 * Hot Rod makes no assumptions about their types.
 *
 * vInt: Variable-length integers are defined defined as compressed,
 * positive integers where the high-order bit of each byte indicates
 * whether more bytes need to be read. The low-order seven bits are
 * appended as increasingly more significant bits in the resulting
 * integer value making it efficient to decode. Hence, values from
 * zero to 127 are stored in a single byte, values from 128 to 16,383
 * are stored in two bytes, and so on.
 *
 * signed vInt: The vInt above is also able to encode negative values,
 * but will always use the maximum size (5 bytes) no matter how small
 * the endoded value is. In order to have a small payload for negative
 * values too, signed vInts uses ZigZag encoding on top of the vInt encoding.
 *
 * vLong : Refers to unsigned variable length long values similar to
 * vInt but applied to longer values. They’re between 1 and 9 bytes long.
 *
 * String : Strings are always represented using UTF-8 encoding.
 ***********************************************************/

/***********************************************************
 * Request Header
 * A request header consists of one of each of the following.
 ***********************************************************/

// Magic (byte)
const reqMagic = 0xA0

// Message ID (vLong)

// Version 2.4 (byte)
const reqVersion = 24

// Opcode (1byte)
const (
	reqOpPut                 = 0x01
	reqOpGet                 = 0x03
	reqOpReplace             = 0x07
	reqOpReplaceIfUnmodified = 0x09
	reqOpRemove              = 0x0B
	reqOpRemoveIfUnmodified  = 0x0D
	reqOpContainsKey         = 0x0F
	reqOpGetWithVersion      = 0x11
	reqOpClear               = 0x13
	reqOpStats               = 0x15
	reqOpPing                = 0x17
	reqOpBulkGet             = 0x19
	reqOpPutAll              = 0x2D
	reqOpGetAll              = 0x2F
)

// Cache Name Length (vInt)

// Cache Name (string)

// Flags (vInt)

// Client Intelligence (byte)
// Basic client, not interested in cluster or hash info.
const reqClientIntel = 0x01

// Topology ID (vInt)
// Always zero for basic clients.
const reqTopologyID = 0

// Transaction Type (byte)
// Initial versions of this client will not support transactions.
const reqTransationType = 0

// Transaction ID (byte array)
// This is omitted when transactions are not supported.
