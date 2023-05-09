"""
Module containing exceptions thrown by :mod:`dxf`.
"""

class DXFError(Exception):
    """
    Base exception class for all dxf errors
    """

class DXFUnexpectedError(DXFError):
    """
    Unexpected value error
    """
    def __init__(self, got, expected):
        """
        :param got: Actual value received
        :param expected: Value that was expected
        """
        super(DXFUnexpectedError, self).__init__()
        self.got = got
        self.expected = expected

class DXFUnexpectedStatusCodeError(DXFUnexpectedError):
    """
    Unexpected HTTP status code
    """
    def __str__(self):
        return 'expected status code %d, got %d' % (self.expected, self.got)

class DXFDigestMismatchError(DXFUnexpectedError):
    """
    Digest didn't match expected value
    """
    def __str__(self):
        return 'expected digest %s, got %s' % (self.expected, self.got)

class DXFUnexpectedKeyTypeError(DXFUnexpectedError):
    """
    Cryptographic key type not supported
    """
    def __str__(self):
        return 'expected key type %s, got %s' % (self.expected, self.got)

class DXFUnexpectedDigestMethodError(DXFUnexpectedError):
    """
    Digest method not supported
    """
    def __str__(self):
        return 'expected digest method %s, got %s' % (self.expected, self.got)

class DXFDisallowedSignatureAlgorithmError(DXFError):
    """
    Signature algorithm forbidden
    """
    def __init__(self, alg):
        """
        :param alg: Forbidden signature algorithm
        :type alg: str
        """
        super(DXFDisallowedSignatureAlgorithmError, self).__init__()
        self.alg = alg

    def __str__(self):
        return 'disallowed signature algorithm: %s' % self.alg

class DXFSignatureChainNotImplementedError(DXFError):
    """
    Signature chains not supported
    """
    def __str__(self):
        return 'verification with a cert chain is not implemented'

class DXFUnauthorizedError(DXFError):
    """
    Registry returned authorized error
    """
    def __str__(self):
        return 'unauthorized'

class DXFAuthInsecureError(DXFError):
    """
    Can't authenticate over insecure (non-HTTPS) connection
    """
    def __str__(self):
        return 'Auth requires HTTPS'

class DXFDigestNotAvailableForSchema1(DXFError):
    """
    https://github.com/docker/distribution/issues/1662#issuecomment-213101772
    A schema1 manifest should always produce the same image id but defining the
    steps to produce directly from the manifest is not straight forward."
    """
    def __str__(self):
        return 'Schema 1 manifest has no configuration blob. See https://github.com/docker/distribution/issues/1662#issuecomment-213101772'

class DXFMountFailed(DXFError):
    """
    Failed to mount blob from another repository.
    """
    def __str__(self):
        return 'Cross repository blob mount failed'

class DXFUnsupportedSchemaType(DXFError):
    """
    Schema type (mediaType) is not recognized/supported
    """
    def __str__(self):
        return 'The mediaType "%s" is not supported' % self.args[0]

class DXFPlatformDataNotFound(DXFError):
    """
    Platform data not found
    """
    def __init__(self, platform):
        """
        :param alg: Platform
        :type alg: str
        """
        super(DXFPlatformDataNotFound, self).__init__()
        self.platform = platform

    def __str__(self):
        return 'Data not found for platform: %s' % self.platform
